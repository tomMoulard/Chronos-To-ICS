package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/julienschmidt/httprouter"
)

var (
	// WeekNumber Number of week to get
	WeekNumber int64
)

// Groups store groups
type Groups []struct {
	Groups []struct {
		Groups []struct {
			Groups []struct {
				Groups []struct {
					Groups []struct {
						Groups   interface{} `json:"Groups"`
						ID       int64       `json:"Id"`
						Name     string      `json:"Name"`
						ParentID int64       `json:"ParentId"`
						Type     int64       `json:"Type"`
					} `json:"Groups"`
					ID       int64  `json:"Id"`
					Name     string `json:"Name"`
					ParentID int64  `json:"ParentId"`
					Type     int64  `json:"Type"`
				} `json:"Groups"`
				ID       int64  `json:"Id"`
				Name     string `json:"Name"`
				ParentID int64  `json:"ParentId"`
				Type     int64  `json:"Type"`
			} `json:"Groups"`
			ID       int64  `json:"Id"`
			Name     string `json:"Name"`
			ParentID int64  `json:"ParentId"`
			Type     int64  `json:"Type"`
		} `json:"Groups"`
		ID       int64  `json:"Id"`
		Name     string `json:"Name"`
		ParentID int64  `json:"ParentId"`
		Type     int64  `json:"Type"`
	} `json:"Groups"`
	ID       int64  `json:"Id"`
	Name     string `json:"Name"`
	ParentID int64  `json:"ParentId"`
	Type     int64  `json:"Type"`
}

// List list of groups, rooms of staff
type List struct {
	ID       int64  `json:"Id"`
	Name     string `json:"Name"`
	ParentID int64  `json:"ParentId"`
	Type     int64  `json:"Type"`
}

// CourseList all informations about a course
type CourseList struct {
	DayList []struct {
		CourseList []struct {
			BeginDate string      `json:"BeginDate"`
			Code      interface{} `json:"Code"`
			Duration  int64       `json:"Duration"`
			EndDate   string      `json:"EndDate"`
			GroupList []List      `json:"GroupList"`
			ID        int64       `json:"Id"`
			Info      string      `json:"Info"`
			Name      string      `json:"Name"`
			RoomList  []List      `json:"RoomList"`
			StaffList []List      `json:"StaffList"`
			Type      string      `json:"Type"`
			URL       string      `json:"Url"`
		} `json:"CourseList"`
		DateTime string `json:"DateTime"`
	} `json:"DayList"`
	ID int64 `json:"Id"`
}

// MergeList extract list name
func MergeList(list []List) string {
	res := ""
	for _, l := range list {
		res += " " + l.Name
	}
	return res
}

// CourseListToIcal convert a CourseList to ical
func CourseListToIcal(cl CourseList, cal *ical.Calendar) {
	for _, dl := range cl.DayList {
		for _, course := range dl.CourseList {
			begin, _ := time.Parse(time.RFC3339, course.BeginDate)
			end, _ := time.Parse(time.RFC3339, course.EndDate)
			description := ""
			if MergeList(course.StaffList) != "" {
				description = "Cours par " + MergeList(course.StaffList) + "\n"
			}
			if MergeList(course.GroupList) != "" {
				description += "Cours pour " + MergeList(course.GroupList) + "\n"
			}
			description += course.Info

			event := cal.AddEvent(fmt.Sprintf("%d-%s-%d@Chronos-to-ICS",
				cl.ID, dl.DateTime, course.ID))
			event.SetCreatedTime(time.Now())
			event.SetDtStampTime(time.Now())
			event.SetModifiedAt(time.Now())
			event.SetStartAt(begin)
			event.SetEndAt(end)
			event.SetSummary(course.Name)
			event.SetDescription(description)
			event.SetLocation(MergeList(course.RoomList))
			event.SetURL(course.URL)
			event.AddAttendee("reciever or participant",
				ical.CalendarUserTypeGroup,
				ical.WithRSVP(false))
			event.SetProperty(ical.ComponentProperty(ical.PropertyDuration),
				fmt.Sprintf("PT%dM", course.Duration))
			for _, grp := range course.GroupList {
				event.AddAttendee(grp.Name + "@epita.fr")
			}
		}
	}
}

// Healthcheck checks is the api is alive
func Healthcheck() int {
	url := "https://v2ssl.webservices.chronos.epita.net/api/v2//Group/GetGroups"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 1
	}

	req.Header.Set("Auth-Token", os.Getenv("ICS_API_KEY"))
	req.Header.Set("Accept", "application/json")

	var client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}

// GetAPI get the value of the API given a query
func GetAPI(query string) (*http.Response, error) {
	url := "https://v2ssl.webservices.chronos.epita.net/api/v2" + query
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Auth-Token", os.Getenv("ICS_API_KEY"))
	req.Header.Set("Accept", "application/json")

	var client = &http.Client{}
	return client.Do(req)
}

// ParseJSON map a json *http.Response into a struct
func ParseJSON(r *http.Response, v interface{}) error {
	if r == nil || r.Body == nil {
		return fmt.Errorf("No Body")
	}

	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// GetFileContent return the content of a file
func GetFileContent(fileName string) (string, error) {
	file, err := os.Open(fileName) // O_RDONLY mode
	if err != nil {
		return "", err
	}
	defer file.Close()

	res, err := ioutil.ReadAll(file)
	return string(res), err
}

// ListGroup Display the group list
func ListGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Printf("Request: %s", "/")
	req, err := GetAPI("/Group/GetGroups")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	var g Groups
	if err = ParseJSON(req, &g); err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fileContent, err := GetFileContent("/html/index.html")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	t, err := template.New("webpage").Parse(fileContent)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	if err := t.Execute(w, g); err != nil {
		log.Println(err)
	}
}

// GetWeek return the weekly calendar for a given entity
func GetWeek(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Printf("Request: /group/%s/%s",
		ps.ByName("groupid"), ps.ByName("entityid"))
	entity := ps.ByName("entityid")
	if entity == "" {
		entity = "1"
	}

	url := "/Week/GetCurrentWeek/" + ps.ByName("groupid") + "/" + entity
	req, err := GetAPI(url)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	var c CourseList
	if err = ParseJSON(req, &c); err != nil {
		fmt.Fprintln(w, err)
		return
	}

	cal := ical.NewCalendar()
	cal.SetMethod(ical.MethodRequest)
	cal.SetName(fmt.Sprintf("Ical for %s", ps.ByName("groupid")))
	cal.SetProductId("-//Tom Moulard//Chronos to ICS")
	CourseListToIcal(c, cal)

	for i := c.ID; i < c.ID+WeekNumber; i++ {
		url = fmt.Sprintf("/Week/GetWeek/%d/%s/%s",
			i, ps.ByName("groupid"), entity)
		req, err = GetAPI(url)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		var currentWeek CourseList
		if err = ParseJSON(req, &currentWeek); err != nil {
			fmt.Fprintln(w, err)
			return
		}
		CourseListToIcal(currentWeek, cal)
	}
	fmt.Fprintln(w, cal.Serialize())
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--healthcheck" {
		os.Exit(Healthcheck())
	}

	IP := os.Getenv("ICS_IP")
	PORT := os.Getenv("ICS_PORT")

	if _, err := fmt.Sscanf(os.Getenv("ICS_WEEK_NUMBER"), "%d", &WeekNumber); err != nil {
		log.Fatalf("'%s' is not a good value for ICS_WEEK_NUMBER",
			os.Getenv("ICS_WEEK_NUMBER"))
	}

	router := httprouter.New()
	router.GET("/", ListGroup)
	router.GET("/group/:groupid/:entityid", GetWeek)
	router.GET("/group/:groupid", GetWeek)

	log.Printf("Starting server at http://%s:%s\n", IP, PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", IP, PORT), router))
}

FROM golang:1.7.3 AS builder

WORKDIR /go/src/github.com/tommoulard/Chronos-To-ISC

RUN go get -d -v \
	github.com/arran4/golang-ical \
	github.com/julienschmidt/httprouter

COPY app.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch
ENV ICS_API_KEY=NOT_PROVIDED
ENV ICS_IP=0.0.0.0
ENV ICS_PORT=8000
ENV ICS_WEEK_NUMBER=4

COPY --from=builder /go/src/github.com/tommoulard/Chronos-To-ISC/app .

COPY ./index.html /html/index.html

USER 1000

HEALTHCHECK  --retries=3 --interval=5s --timeout=3s \
	CMD ./app --healthcheck

CMD ["./app"]

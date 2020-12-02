# Chronos-To-ICS

See [Chronos](http://chronos.epita.net/), [Chronos API](https://v2ssl.webservices.chronos.epita.net/api/v2).

## Setup
```bash
git clone https://github.com/tommoulard/Chronos-To-ICS && cd Chronos-To-ICS
cp .env.default .env
```

And change values in the `.env` file

## Usage
```bash
make
```

### In another docker-compose
```yml
version: '3.6'

services:
  ics:
    image: tommoulard/chronos-to-ics
    ports:
      - "8000:8000"
    environment:
      - 'ICS_API_KEY=${API_KEY}'
      - 'ICS_PORT=8000'
      - 'ICS_WEEK_NUMBER=4'
    volumes:
      - '/etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt'
```

Where:

 - `ICS_API_KEY` is the Chronos API key
 - `ICS_PORT` is the port on which the application must bind, it should be the same as the one exposed with docker.
 - `ICS_WEEK_NUMBER` is the number of week to fetch, the higher the more event fetched, but the more API calls for each calendars. `4` week is `1` month.

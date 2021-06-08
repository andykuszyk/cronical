# cronical
Cronical is a web server capable of filtering out events from a WebCal feed based on a Cron expression. The idea behind this is that you may have a calendar feed which you want to import into a calendar application, but you want to exclude events from a certain time range.

My motivtion for creating it, is that I want to add my on-call calendar to my personal calendar, but I don't want to include on-call shifts during working hours--just the ones in the evenings/weekends that affect my personal life.

## Usage
The basic usage of Cronical (assuming it's running on port 8080) is:

```sh
curl http://localhost:8080/filter?ical=<ical>&exclude=<cron>
```

Where the `<ical>` parameter is the base64 encoded URL of the WebCal to filter, and the `<cron>` parameter is the base64 encoded CRON expression to use when filtering out events.

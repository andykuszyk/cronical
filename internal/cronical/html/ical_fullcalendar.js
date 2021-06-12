// Depends on ./ical_events.js

recur_events = []

function an_filter(string) {
    // remove non alphanumeric chars
    return string.replace(/[^\w\s]/gi, '')
}

function moment_icaltime(moment, timezone) {
    // TODO timezone
    return new ICAL.Time().fromJSDate(moment.toDate())
}

function expand_recur_events(start, end, timezone, events_callback) {
    events = []
    for (event of recur_events) {
        expand_recur_event(event, moment_icaltime(start, timezone), moment_icaltime(end, timezone), function(event){
            fc_event(event, function(event){
                events.push(event)
            }, ['recur-event'])
        })
    }
    events_callback(events)
}

function fc_events(ics) {
    events = []
    ical_events(ics,
        function(event){
            fc_event(event, function(event){
                events.push(event)
            })
        },
        function(event){
            recur_events.push(event)      
        })
    return events
}

function fc_event(event, event_callback, class_list) {
    e = {
        title:event.getFirstPropertyValue('summary'),
        url:event.getFirstPropertyValue('url'),
        id:event.getFirstPropertyValue('uid'),
        className:['event-'+an_filter(event.getFirstPropertyValue('uid'))].concat(class_list),
        allDay:false
    }
    try {
        e['start'] = event.getFirstPropertyValue('dtstart').toJSDate()
    } catch (TypeError) {
        console.debug('Undefined "dtstart", vevent skipped.')
        return
    }
    try {
        e['end'] = event.getFirstPropertyValue('dtend').toJSDate()
    } catch (TypeError) {
        e['allDay'] = true
    }
    event_callback(e)
}


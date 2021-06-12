ics_sources = ['http://localhost:8080/webcal?ical=aHR0cHM6Ly9mb3JtMy5wYWdlcmR1dHkuY29tL3ByaXZhdGUvZjBiYzBjNGY0N2ViYWY5OWM2MDNiY2NiN2IwYjIyMTUyMzNmZjcwY2']


function data_req (url, callback) {
    req = new XMLHttpRequest()
    req.addEventListener('load', callback)
    req.open('GET', url)
    req.send()
}

function add_recur_events() {
    if (sources_to_load_cnt < 1) {
        $('#calendar').fullCalendar('addEventSource', expand_recur_events)
    } else {
        setTimeout(add_recur_events, 30)
    }
}

$(document).ready(function() {
    $('#calendar').fullCalendar({
        header: {
            left: 'prev,next today',
            center: 'title',
            right: 'month,agendaWeek,agendaDay'
        },
        defaultView: 'month',
	defaultDate: '2016-03-01'
    })
    sources_to_load_cnt = ics_sources.length
    for (ics of ics_sources) {
        data_req(ics, function(){
            $('#calendar').fullCalendar('addEventSource', fc_events(this.response))
            sources_to_load_cnt -= 1
        })
    }
    add_recur_events()
})


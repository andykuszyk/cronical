ics = 'http://localhost:8080/filter?ical=aHR0cHM6Ly9mb3JtMy5wYWdlcmR1dHkuY29tL3ByaXZhdGUvZjBiYzBjNGY0N2ViYWY5OWM2MDNiY2NiN2IwYjIyMTUyMzNmZjcwY2FmNzM0ODY3OWZmYTBjMDFlNGZkM2M2YS9mZWVk&exclude=KiA5LTE2ICogKiAq'

function generate() {
    var webcalUrl = $( '#webcal_url' ).val()
    var cronExpression = $( '#cron' ).val()
    console.log('generate clicked with webcal: ' + webcalUrl + ' and cron: ' + cronExpression)
    webcalUrlBase64 = btoa(webcalUrl)
    cronExpressionBase64 = btoa(cronExpression)
    var generatedUrl = location.protocol + '//' + location.host + '/filter?ical=' + webcalUrlBase64 + '&exclude=' + cronExpressionBase64
    $( '#filtered_url' ).val(generatedUrl)

    data_req(generatedUrl, function(){
        $('#calendar').fullCalendar('addEventSource', fc_events(this.response))
    })
}

function data_req (url, callback) {
    req = new XMLHttpRequest()
    req.addEventListener('load', callback)
    req.open('GET', url)
    req.send()
}

$(document).ready(function() {
    var d = new Date()
    $('#calendar').fullCalendar({
        header: {
            left: 'prev,next today',
            center: 'title',
            right: 'month,agendaWeek,agendaDay'
        },
        defaultView: 'month',
        defaultDate: d.getFullYear() + '-' + d.getMonth() + '-' + d.getDate()
    })
})


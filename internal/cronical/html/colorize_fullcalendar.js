function hash_code(string) {
    h = 0
    for (i = 0; i < string.length; i++) {
        h = string.charCodeAt(i) + ((h << 5) - h);
    }
    return h
}

function hash_color(string) {
    return 'hsl('+hash_code(string)%360+', 70%, 50%)'
}

function colorize_event(e) {
    for (c of e.classList) {
       if (c.startsWith('event-')) {
           e.style['background'] = e.style['border-color'] = hash_color(c)
       }
    }     
}

function colorize_events() {
    ee = document.getElementsByClassName('fc-event')
    for (i in ee) {
        try {
            colorize_event(ee[i])
        } catch (TypeError) {}
    }
}



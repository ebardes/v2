var ws

function openSocket() {
  var dest = "ws://" + location.host + "/ws/"
  console.log(document.URL)
  console.log(dest)
  ws = new WebSocket(dest);

  var func = {
    "ack": ackFunction,
    "pkt": packetFunction,
  }

  ws.onopen = function(evt) {
    console.log("OPEN");
    ws.send(JSON.stringify({verb: "reg", display: display}))
  }
  ws.onclose = function(evt) {
    console.log("CLOSE");
    ws = null;
  }
  ws.onmessage = function(evt) {
    obj = JSON.parse(evt.data)
    f = func[obj.verb]
    if (f) {
      f(obj)
    } else {
      console.log("UNKNOWN RESPONSE: " + evt.data);
    }
  }
  ws.onerror = function(evt) {
    console.log("ERROR: " + evt.data);
  }
}

function ackFunction(e) {
  console.log(e)
}

function packetFunction(e) {
  console.log(
    "I got a packet!",
    e
  )
}

$(document).ready(function(){
  openSocket()
})


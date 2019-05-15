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
  setupViewport(e)
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

var canvasmap = {}

function setupViewport(cfg) {
  $('#viewport').html('')
  var len = cfg.layers.length
  for (var i = 0; i < len; i++) {
    var layerid = "layer_"+cfg.layers[i]
    var x = $('#viewport').append('<div id="'+layerid+'"><canvas id="c'+layerid+'"></canvas><img id="i'+layerid+'"/>')
    var c = document.getElementById("c"+layerid)
    canvasmap[cfg.layers[i]] = c
    c.width = 1920
    c.height = 1080
  }
  ws.send(JSON.stringify({verb: "ref", display: display}))
}
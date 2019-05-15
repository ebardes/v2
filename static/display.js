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
  setupViewport(e)
}

var sizex = 1920
var sizey = 1080
var centerx = 0
var centery = 0

function packetFunction(e) {
  var layer = e.layer
  var packet = e.packet
  // if (packet.Changed) {
    $('#ilayer_'+layer).attr('src', packet.URL)
  // }

  var img = document.getElementById("ilayer_"+layer)
  var c = document.getElementById("clayer_"+layer)
  
  $(c).attr('style', 'filter: opacity(' + (packet.Level / 255) + ');')
  var gl1 = c.getContext("2d")

  var composite = document.createElement("canvas")
  composite.width = sizex
  composite.height = sizey

  var gl = composite.getContext("2d")
  gl.fillStyle='#000000'
  gl.fillRect(0,0,sizex, sizey)
  gl.translate(centerx, centery)

  // bright = (packet.Brightness / 127.0) + 1.0
  // gl.filter = "brightness("+bright+")"
  // packet.bright = bright

  gl.rotate((packet.ZRotate) / 32768.0)
  gl.drawImage(img, -centerx, -centery, img.width, img.height)
  
  gl1.drawImage(composite, 0, 0)

  if e.layer == 1 {
    console.log(packet)
  }

  
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
    c.width = sizex
    c.height = sizey
  }
  ws.send(JSON.stringify({verb: "ref", display: display}))
}
<html>
<head>
  <base href="{{.Prefix}}"/>
  <link rel="stylesheet" href="{{.Prefix}}display.css" type="text/css" />
  <script type="text/javascript" src="{{.Prefix}}jquery-3.2.1.min.js"></script>
  <script type="text/javascript" src="{{.Prefix}}display.js"></script>
  <script type="text/javascript">var panel={{.Data.Panel}};</script>
</head>
<body>
  {{ range $i,$v :=  .Data.Layers }}
  <div id="layer_{{$i}}">{{$v}}</div>
  {{ end }}
</body>
</html>

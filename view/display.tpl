<html>
<head>
  <base href="{{.Prefix}}"/>
  <link rel="stylesheet" href="{{.Prefix}}display.css" type="text/css" />
  <script type="text/javascript" src="{{.Prefix}}jquery-3.2.1.min.js"></script>
  <script type="text/javascript" src="{{.Prefix}}display.js"></script>
  <script type="text/javascript">var display={{- .Data.ID -}};</script>
</head>
<body>
{{- range $i,$x := .Data.Layers -}}
  <canvas id="main_{{$i}}" class="main"></canvas>
{{- end -}}
</body>
</html>

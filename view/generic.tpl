{{define "head"}}
  <html>
  <head>
    <link rel="icon" type="image/png" href="/server.png" />
  	<link rel="stylesheet" href="{{.Prefix}}style.css" type="text/css" />
  	<script type="text/javascript" src="{{.Prefix}}jquery-3.2.1.min.js"></script>
  	<script type="text/javascript" src="{{.Prefix}}main.js"></script>
    <title>{{.Title}}</title>
  </head>
  <body>
{{end}}
{{define "foot"}}
  </body>
</html>
{{end}}

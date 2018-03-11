{{template "head"}}
	<div class="interactive">
		<h1>Network</h1>
		<div>Universe: <b>{{.Data.Universe}}</b></div>
	</div>

	<div class="interactive">
		{{- range $a := .Panes -}}
		<h1>Pane {{.ID}}</h1>
		<div>Starting Address: <b>{{.Data.StartAddress}}</b></div>
		{{- end -}}
	</div>
{{template "foot"}}

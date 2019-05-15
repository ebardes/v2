{{- template "head" -}}
	<div class="interactive">
		<h1>Network</h1>
		<div>Universe: <b>{{.Data.Universe}}</b></div>
		<div>Protocol: <b>{{.Data.Protocol}}</b></div>
		<div>Adapter: <b>{{.Data.Interface}}</b></div>
	</div>

	<div class="interactive">
		{{- range $a := .Data.Displays -}}
		<h1>Display {{.ID}}</h1>
			{{- range $i,$a := .Layers -}}
			<h2>Layer {{$i}}</h2>
			<div>Starting Address: <b>{{.StartAddress}}</b></div>
			{{- end -}}
		{{- end -}}
	</div>

	<div class="interactive">
		<h1>Content</h1>
		{{- range $n, $group := .Data.Content -}}
		<h3>Group {{$n}}</h3>
			{{- range $id, $slot := .Slots -}}
			<div class="slot">
				<div class="head">{{- $id -}}: {{ .Name -}}</div>
				<div class="body"><img src="{{- .URL }}" /></div>
			</div>
			{{- end -}}
			{{- if ne 0 $n -}}
			<div class="slot">
				<div class="head">New Content</div>
				<div class="new body" group="{{$n}}">
				</div>
			</div>
			{{- end -}}
		{{- end -}}
	</div>

	<pre>{{.Data}}</pre>
{{- template "foot" -}}

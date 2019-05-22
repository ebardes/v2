{{- template "head" -}}
	<div class="interactive">
		<h1>Network</h1>
		<form id="netcfg">
		<div>Universe: <input name="universe" value="{{.Data.Universe}}"></div>
		{{- $p := .Data.Protocol -}}
		<div>Protocol:
			<select name="protocol">
			{{- range .Data.Protocols -}}
			<option {{ if eq . $p -}}selected="selected" {{ end -}}value="{{.}}">{{.}}</option>
			{{- end -}}
			</select>
		</div>
		{{- $n := .Data.Interface -}}
		<div>Adapter:
			<select name="network">
			{{- range .Data.Networks -}}
			<option {{ if eq .Name $n -}}selected="selected" {{ end -}}value="{{.Name}}">{{.Name}} - {{.IPAddress}}</option>
			{{- end -}}
			</select>
		</div>
		</form>
		<input type="button" value="Change" />
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
				<div class="head"><b>{{- $id -}}</b>: {{ .Name -}}</div>
				<div class="body"><img src="{{- .URL }}" />
					{{- if ne 0 $n -}}	
					<div class="delete" id="/delete/{{$n}}/{{$id}}">Del</div>
					{{- end -}}
				</div>
			</div>
			{{- end -}}
			{{- if ne 0 $n -}}
			<div class="slot">
				<div class="head">New Content</div>
				<div class="new body" group="{{$n}}"></div>
			</div>
			{{- end -}}
		{{- end -}}
	</div>
{{- template "foot" -}}

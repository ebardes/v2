package view

import (
	"html/template"
	"io"
	"log"
	"os"
)

type View struct {
	w    io.Writer
	info Info
}

type Info struct {
	Prefix string
	Title  string
	Data   interface{}
}

// AllTemplates is the root template
var AllTemplates *template.Template

func init() {
	gopath := os.Getenv("GOPATH")

	g, err := template.ParseGlob(gopath + "/view/*.tpl")
	if err != nil {
		log.Panic(err)
	}

	AllTemplates = g
}

func InitView(writer io.Writer) *View {
	info := Info{}
	return &View{w: writer, info: info}
}

func (me *View) Render(x string, data interface{}) error {
	me.info.Data = data
	return AllTemplates.ExecuteTemplate(me.w, x, me.info)
}

func (me *View) Title(title string) *View {
	me.info.Title = title
	return me
}

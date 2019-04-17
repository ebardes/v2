package view

import (
	"html/template"
	"io"
	"path"
	"v2/config"

	"github.com/rs/zerolog/log"
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

func Init(cfg *config.Config) {
	p := path.Join(cfg.TemplateDir, "*.tpl")
	g, err := template.ParseGlob(p)
	if err == nil {
		AllTemplates = g
	} else {
		log.Error().Err(err).Msg("Could not templates")
	}
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

func (me *View) Prefix(prefix string) *View {
	me.info.Prefix = prefix
	return me
}

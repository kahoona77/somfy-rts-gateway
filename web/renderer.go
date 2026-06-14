package web

import (
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"

	"somfyRtsGateway/core"
)

type Template struct {
	templates *template.Template
}

func NewTemplate(ctx *core.Ctx) *Template {
	templates := template.Must(template.New("base.html").Funcs(template.FuncMap{
		"basePath": func() string {
			return fmt.Sprintf("%s/", ctx.Config.BasePath)
		},
	}).ParseFS(mustSub("tmpl"), "*.html"))

	return &Template{templates: templates}
}

func (t *Template) Render(w io.Writer, _ string, data interface{}, _ echo.Context) error {
	return t.templates.ExecuteTemplate(w, "base.html", data)
}

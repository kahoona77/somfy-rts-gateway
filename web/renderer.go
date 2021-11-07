package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"somfyRtsGateway/core"
)

type Template struct {
	ctx *core.Ctx
}

func NewTemplate(ctx *core.Ctx) *Template {
	return &Template{ctx: ctx}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	temp, err := template.New("base.html").Funcs(template.FuncMap{
		"basePath": func() string {
			return fmt.Sprintf("%s/", t.ctx.Config.BasePath)
		},
	}).ParseFiles("web/tmpl/base.html", "web/tmpl/"+name)
	if err != nil {
		logrus.Errorf("error parsing templates: %v", err)
		return err
	}
	return temp.ExecuteTemplate(w, temp.Name(), data)
}

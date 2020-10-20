package core

import (
	"github.com/labstack/echo/v4"
	"strconv"
)

func CreateCtx(ctx *Ctx) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &WebContext{Context: c, Ctx: ctx}
			return next(cc)
		}
	}
}

type Ctx struct {
	AppConfig *AppConfig
}

func (ctx *Ctx) Copy() *Ctx {
	return &Ctx{}
}

func (ctx *Ctx) Close() {
}

type WebContext struct {
	echo.Context
	*Ctx
}

func (ctx *WebContext) ParamAsInt(name string) int {
	p, _ := strconv.Atoi(ctx.Param(name))
	return p
}

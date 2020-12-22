package views

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"somfyRtsGateway/core"
)

func Index(c echo.Context) error {
	ctx := c.(*core.WebContext).Ctx

	data := struct {
		Devices []core.Device
	}{Devices: ctx.Controller.Devices()}

	return c.Render(http.StatusOK, "index.html", &data)
}

func Cmd(c echo.Context) error {
	ctx := c.(*core.WebContext).Ctx

	device := c.Param("device")
	cmd := c.Param("cmd")

	ctx.CommandChannel <- core.DeviceCmd{
		Device: device,
		Cmd:    cmd,
	}

	return c.Redirect(http.StatusFound, ctx.Config.AbsolutePath("web"))
}

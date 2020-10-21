package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"somfyRtsGateway/core"
)

func SomfyCmd(ec echo.Context) error {
	ctx := ec.(*core.WebContext)

	device := ctx.Param("device")
	cmd := ctx.Param("cmd")

	ctx.CommandChannel <- core.DeviceCmd{
		Device: device,
		Cmd:    cmd,
	}

	return ctx.String(http.StatusOK, fmt.Sprintf("issued command %s to device %s", cmd, device))
}

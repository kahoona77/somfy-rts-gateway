package somfy

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"somfyRtsGateway/core"
)

func Cmd(ec echo.Context) error {
	ctx := ec.(*core.WebContext)

	device := ctx.Param("device")
	cmd := ctx.Param("cmd")

	ctx.CommandChannel <- core.DeviceCmd{
		Device: device,
		Cmd:    cmd,
	}

	return ctx.String(http.StatusOK, fmt.Sprintf("issued command %s to device %s", cmd, device))
}

func ListDevices(ctrl *Controller) func(echo.Context) error {
	return func(ec echo.Context) error {
		return ec.JSON(http.StatusOK, ctrl.GetDevices())
	}
}

func GetDevice(ctrl *Controller) func(echo.Context) error {
	return func(ec echo.Context) error {
		device := ec.Param("device")
		d, err := ctrl.GetDevice(device)
		if err != nil {
			return ec.String(http.StatusNotFound, err.Error())
		}
		return ec.JSON(http.StatusOK, d)
	}
}

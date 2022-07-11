package controller

import (
	"fmt"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"github.com/gin-gonic/gin"
)

var DefaultTimeout = 30
var systemCtlOpts = systemctl.Options{
	UserMode: false,
	Timeout:  DefaultTimeout,
}

type SystemCtl struct {
	Action      string `json:"action" validate:"required"` // start, stop, restart
	ServiceName string `json:"service_name" validate:"required"`
}

func getSystemCtlBody(c *gin.Context) (dto *SystemCtl, err error) {
	err = c.ShouldBindJSON(&dto)
	return dto, err
}

func (inst *Controller) SystemCtl(c *gin.Context) {
	body, err := getSystemCtlBody(c)
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	out, err := systemctl.Status(body.ServiceName, systemCtlOpts)
	fmt.Println("out>>>", out)
	fmt.Println("err>>>", err)
	fmt.Println("body.ServiceName>>>", body.ServiceName)
}

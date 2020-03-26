package app

import (
	"github.com/gin-gonic/gin"
	"github.com/onethefour/REST-GO-demo/app/controller"
)

func Router(r *gin.Engine) {
	//new(controller.TestController).Router(r)

	new(controller.AccountController).Router(r)
	new(controller.StaticController).Router(r)
	new(controller.LianghuaController).Router(r)
	//启动服务
	//go monitor.NewMonitor().Start()
}

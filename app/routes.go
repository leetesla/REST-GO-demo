package app

import (
	"github.com/gin-gonic/gin"
	"REST-GO-demo/app/controller"
)

func Router(r *gin.Engine) {

	new(controller.AccountController).Router(r)
	new(controller.LianghuaController).Router(r)
	new(controller.StaticController).Router(r)

}

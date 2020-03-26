package main

import (
	"log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/onethefour/REST-GO-demo/app/utils"
	"github.com/onethefour/REST-GO-demo/app"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}
func main() {

	//gin.DefaultWriter = utils.Loger()
	//gin.SetMode(gin.DebugMode)
	//log.SetOutput(gin.DefaultWriter)
	router := gin.Default()
	//session中间件
	router.Use(sessions.Sessions("mysession", utils.Store))
	//login_token校验登陆中间件
	router.Use(utils.CheckLogin)

	app.Router(router)

	router.Run(":8000")

}

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StaticController struct{}

func (ctr *StaticController) Router(r *gin.Engine) {
	r.StaticFS("/front", http.Dir("./element-html/dist"))
}

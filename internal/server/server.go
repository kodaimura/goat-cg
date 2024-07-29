package server

import (
	"github.com/gin-gonic/gin"

	"goat-cg/config"
)

func Run() {
	cf := config.GetConfig()
	r := router()
	r.Run(":" + cf.AppPort)
}

func router() *gin.Engine {
	r := gin.Default()
	
	//TEMPLATE
	r.LoadHTMLGlob("web/template/*.html")

	//STATIC
	r.Static("/css", "web/static/css")
	r.Static("/js", "web/static/js")
	r.Static("/tmp", "./tmp")
	r.StaticFile("/favicon.ico", "web/static/favicon.ico")

	SetRouter(r)

	return r
}

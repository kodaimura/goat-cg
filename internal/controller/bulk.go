package controller

import (
	"github.com/gin-gonic/gin"

	"goat-cg/internal/shared/constant"
	"goat-cg/internal/service"
)


type bulkController struct {
	tServ service.TableService
	cServ service.ColumnService
	urlServ service.UrlCheckService
}


func newBulkController() *bulkController {
	tServ := service.NewTableService()
	cServ := service.NewColumnService()
	urlServ := service.NewUrlCheckService()
	return &bulkController{tServ, cServ, urlServ}
}


//GET /:project_cd/bulk-registration-ddl
func (ctr *bulkController) bulkPage(c *gin.Context) {
	ctr.urlServ.CheckProjectCdAndGetProjectId(c)

	c.HTML(200, "bulk.html", gin.H{
		"commons": constant.Commons,
		"project_cd": c.Param("project_cd"), 
	})
}
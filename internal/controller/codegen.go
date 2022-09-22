package controller

import (
	"github.com/gin-gonic/gin"

	"goat-cg/pkg/utils"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/service"
)


type codegenController struct {
	tServ service.TableService
	cgServ service.CodegenService
	urlServ service.UrlCheckService
}


func newCodegenController() *codegenController {
	tServ := service.NewTableService()
	cgServ := service.NewCodegenService()
	urlServ := service.NewUrlCheckService()
	return &codegenController{tServ, cgServ, urlServ}
}


//GET /:project_cd/codegen
func (ctr *codegenController) codegenPage(c *gin.Context) {
	projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)
	
	tables, _ := ctr.tServ.GetTables(projectId)

	c.HTML(200, "codegen.html", gin.H{
		"commons": constant.Commons,
		"project_cd": c.Param("project_cd"), 
		"tables": tables,
	})
}


type CodegenPostBody struct {
	TableIds []string `json:"tableids"`
	DbType string `json:"dbtype"`
}


//POST /:project_cd/codegen/goat
func (ctr *codegenController) codegenGOAT(c *gin.Context) {
	ctr.urlServ.CheckProjectCdAndGetProjectId(c)

	pb := &CodegenPostBody{} 
	c.BindJSON(&pb)

	tableIds, err := utils.AtoiSlice(pb.TableIds)

	if err != nil {
		c.String(200, "error.txt")
		return
	}

	fpath := ctr.cgServ.CodeGenerateGoat(pb.DbType, tableIds)

	c.String(200, fpath[1:])
}


//POST /:project_cd/codegen/ddl
func (ctr *codegenController) codegenDDL(c *gin.Context) {
	ctr.urlServ.CheckProjectCdAndGetProjectId(c)

	pb := &CodegenPostBody{} 
	c.BindJSON(&pb)

	tableIds, err := utils.AtoiSlice(pb.TableIds)

	if err != nil {
		c.String(200, "error.txt")
		return
	}

	fpath := ctr.cgServ.CodeGenerateDdl(pb.DbType, tableIds)

	c.String(200, fpath[1:])
}

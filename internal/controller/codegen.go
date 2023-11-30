package controller

import (
	"github.com/gin-gonic/gin"

	"goat-cg/pkg/utils"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/service"
)


type CodegenController struct {
	tableService service.TableService
	codegenService service.CodegenService
}


func NewCodegenController() *CodegenController {
	tableService := service.NewTableService()
	codegenService := service.NewCodegenService()
	return &CodegenController{tableService, codegenService}
}


//GET /:username/:project_name/codegen
func (ctr *CodegenController) CodegenPage(c *gin.Context) {
	p := c.Keys["project"].(model.Project)
	
	tables, _ := ctr.tableService.GetTables(p.ProjectId)

	c.HTML(200, "codegen.html", gin.H{
		"commons": constant.Commons,
		"project_name": c.Param("project_name"), 
		"tables": tables,
	})
}


type CodegenPostBody struct {
	TableIds []string `json:"tableids"`
	DbType string `json:"dbtype"`
}


//POST /:username/:project_name/codegen/goat
func (ctr *CodegenController) CodegenGOAT(c *gin.Context) {
	pb := &CodegenPostBody{} 
	c.BindJSON(&pb)

	tableIds, err := utils.AtoiSlice(pb.TableIds)

	if err != nil {
		c.String(200, "error.txt")
		return
	}

	fpath := ctr.codegenService.CodeGenerateGoat(pb.DbType, tableIds)

	c.String(200, fpath[1:])
}


//POST /:project_cd/codegen/ddl
func (ctr *CodegenController) CodegenDDL(c *gin.Context) {
	pb := &CodegenPostBody{} 
	c.BindJSON(&pb)

	tableIds, err := utils.AtoiSlice(pb.TableIds)

	if err != nil {
		c.String(200, "error.txt")
		return
	}

	fpath := ctr.codegenService.CodeGenerateDdl(pb.DbType, tableIds)

	c.String(200, fpath[1:])
}

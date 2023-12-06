package controller

import (
	"github.com/gin-gonic/gin"

	"goat-cg/pkg/utils"
	"goat-cg/internal/service"
	"goat-cg/internal/model"
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
func (cc *CodegenController) CodegenPage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)
	
	tables, _ := cc.tableService.GetTables(project.ProjectId)

	c.HTML(200, "codegen.html", gin.H{
		"project": project,
		"tables": tables,
	})
}


type CodegenPostBody struct {
	TableIds []string `json:"tableids"`
	DbType string `json:"dbtype"`
}


//POST /:username/:project_name/codegen/goat
func (cc *CodegenController) CodegenGOAT(c *gin.Context) {
	pb := &CodegenPostBody{} 
	c.BindJSON(&pb)

	tableIds, err := utils.AtoiSlice(pb.TableIds)

	if err != nil {
		c.String(200, "error.txt")
		return
	}

	fpath := cc.codegenService.CodeGenerateGoat(pb.DbType, tableIds)

	c.String(200, fpath[1:])
}


//POST /:project_cd/codegen/ddl
func (cc *CodegenController) CodegenDDL(c *gin.Context) {
	pb := &CodegenPostBody{} 
	c.BindJSON(&pb)

	tableIds, err := utils.AtoiSlice(pb.TableIds)

	if err != nil {
		c.String(200, "error.txt")
		return
	}

	fpath := cc.codegenService.CodeGenerateDdl(pb.DbType, tableIds)

	c.String(200, fpath[1:])
}

package controller

import (
    "github.com/gin-gonic/gin"

    "goat-cg/internal/core/jwt"
    "goat-cg/internal/constant"
    "goat-cg/internal/service"
)


type projectController struct {
    pServ service.ProjectService
}


func newProjectController() *projectController {
    pServ := service.NewProjectService()
    return &projectController{pServ}
}

//GET /projects
func (ctr *projectController) projectsPage(c *gin.Context) {
    projects, _ := ctr.pServ.GetProjects(jwt.GetUserId(c))
    
    c.HTML(200, "projects.html", gin.H{
        "commons": constant.Commons,
        "projects": projects,
    })
}


//GET /projects/new
func (ctr *projectController) projectPage(c *gin.Context) {
    
    c.HTML(200, "project.html", gin.H{
        "commons": constant.Commons,
    })
}


//POST /projects
func (ctr *projectController) postProjects(c *gin.Context) {
    result := ctr.pServ.CreateProject(
        jwt.GetUserId(c), 
        c.PostForm("project_cd"), 
        c.PostForm("project_name"),
    )
    
    if result == service.CREATE_PROJECT_SUCCESS_INT {
        c.Redirect(303, "/projects")

    } else if result == service.CREATE_PROJECT_CONFLICT_INT {
        c.HTML(409, "project.html", gin.H{
            "commons": constant.Commons,
            "error": "ProjectCd が既に使われています。",
        })

    } else {
        c.HTML(500, "project.html", gin.H{
            "commons": constant.Commons,
            "error": "登録に失敗しました。",
        })

    }
}


package controller

import (
    "strconv"
    "github.com/gin-gonic/gin"

    "goat-cg/internal/core/jwt"
    "goat-cg/internal/shared/constant"
    "goat-cg/internal/service"
    "goat-cg/internal/model/entity"
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
    userId := jwt.GetUserId(c)
    projectCd := c.Query("project_cd")
    var project entity.Project

    projects, _ := ctr.pServ.GetProjects(userId)
    projects2, _ := ctr.pServ.GetProjectsPendingApproval(userId)
    
    if projectCd != "" {
        project = ctr.pServ.GetProjectByCd(projectCd)
    }

    c.HTML(200, "projects.html", gin.H{
        "commons": constant.Commons,
        "projects": projects,
        "projects2": projects2,
        "project":project,
    })
}


//GET /projects/:project_id/join
func (ctr *projectController) joinRequest(c *gin.Context) {
    userId := jwt.GetUserId(c)
    projectId, err := strconv.Atoi(c.Param("project_id"))

    if err != nil {
        c.Redirect(303, "/projects")
    }

    ctr.pServ.JoinRequest(userId, projectId)
    
    c.Redirect(303, "/projects")
}


//GET /projects/:project_id/cancel
func (ctr *projectController) cancelJoinRequest(c *gin.Context) {
    userId := jwt.GetUserId(c)
    projectId, err := strconv.Atoi(c.Param("project_id"))

    if err != nil {
        c.Redirect(303, "/projects")
    }

    ctr.pServ.CancelJoinRequest(userId, projectId)
    
    c.Redirect(303, "/projects")
}


//GET /projects/new
func (ctr *projectController) createProjectPage(c *gin.Context) {
    
    c.HTML(200, "project.html", gin.H{
        "commons": constant.Commons,
    })
}


//POST /projects
func (ctr *projectController) createProject(c *gin.Context) {
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


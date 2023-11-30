package middleware

import (
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
)


func PathParameterValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := jwt.GetUsername(c)
		ownername := c.Param("username")
		projectName := c.Param("project_name")
		tableId := c.Param("tableId")

		p, _ := searchAccessibleProject(username, ownername, projectName) 
		
		if p.projectId == 0 {
			c.HTML(404, "404error.html", gin.H{})
			c.Abort()
			return
		} 

		c.Set("project", p)

		if tableId != "" {
			t, _ := searchAccessibleTable(p.ProjectId, tableId)

			if t.TableId == 0 {
				c.HTML(404, "404error.html", gin.H{})
				c.Abort()
				return
			}

			c.Set("table", t)
		}

		c.Next()
	}
}

func searchAccessibleProject (username, ownername, projectName string) (model.Project, error) {
	pr := repository.NewProjectService()

	if username == ownername {
		return pr.GetByUniqueKey(username, projectName)
	} else {
		return pr.GetMemberProject(username, projectName)
	}
}

func searchAccessibleTable (projectId, tableId string) (model.Table, error) {
	tr := repository.NewTableRepository()
	return tr.GetByIdAndProjectId(tableId, projectId)
}
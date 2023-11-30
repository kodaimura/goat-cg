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
		tableId := c.Param("table_id")
		columnId := c.Param("column_id")

		project, b := searchAccessibleProject(username, ownername, projectName) 
		
		if !b {
			c.HTML(404, "404error.html", gin.H{})
			c.Abort()
			return
		} 

		c.Set("project", project)

		if tableId != "" {
			table, b := searchAccessibleTable(project.ProjectId, tableId)

			if b! {
				c.HTML(404, "404error.html", gin.H{})
				c.Abort()
				return
			}

			c.Set("table", table)
		}

		if columnId != "" {
			column, b := searchAccessibleColumn(tableId, columnId)

			if b! {
				c.HTML(404, "404error.html", gin.H{})
				c.Abort()
				return
			}

			c.Set("column", column)
		}

		c.Next()
	}
}

func searchAccessibleProject (username, ownername, projectName string) (model.Project, bool) {
	var p model.Project
	pr := repository.NewProjectService()

	if username == ownername {
		p, _ = pr.GetByUniqueKey(username, projectName)
	} else {
		p, _ = pr.GetMemberProject(username, projectName)
	}

	if p.projectId == 0 {
		return p, false
	}
	return p, true
}

func searchAccessibleTable (projectId, tableId string) (model.Table, bool) {
	tr := repository.NewTableRepository()
	t, _ = tr.GetById(tableId)

	if t.ProjectId != projectId {
		return t, false
	}
	return t, true
}

func searchAccessibleColumn (projectId, tableId string) (model.Column, bool) {
	cr := repository.NewColumnRepository()
	c, _ = cr.GetById(columnId)
	
	if c.TableId != tableId {
		return c, false
	}
	return c, true
}
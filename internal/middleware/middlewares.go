package middleware

import (
	"errors"
	"strconv"
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
	"goat-cg/internal/query"
)


func PathParameterValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		project, err := validateProjectNameAndGetProject(c)
		if err != nil {
			c.HTML(404, "404error.html", gin.H{})
			c.Abort()
			return
		} 

		c.Set("project", project)

		if c.Param("table_id") != "" {
			tableId, err := strconv.Atoi(c.Param("table_id"))
			if err != nil {
				c.HTML(404, "404error.html", gin.H{})
				c.Abort()
				return
			}

			table, err := validateTableIdAndGetTable(project.ProjectId, tableId)
			if err != nil {
				c.HTML(404, "404error.html", gin.H{})
				c.Abort()
				return
			}

			c.Set("table", table)

			if c.Param("column_id") != "" {
				columnId, err := strconv.Atoi(c.Param("column_id"))
				if err != nil {
					c.HTML(404, "404error.html", gin.H{})
					c.Abort()
					return
				}
	
				column, err := validateColumnIdAndGetColumn(table.TableId, columnId)
				if err != nil {
					c.HTML(404, "404error.html", gin.H{})
					c.Abort()
					return
				}
	
				c.Set("column", column)
			}
		}

		c.Next()
	}
}

func validateProjectNameAndGetProject (c *gin.Context) (model.Project, error) {
	userId := jwt.GetUserId(c)
	username := jwt.GetUsername(c)
	ownername := c.Param("username")
	projectName := c.Param("project_name")

	var err error
	var p model.Project

	if username == ownername {
		pr := repository.NewProjectRepository()
		p, err = pr.GetOne(&model.Project{Username: username, ProjectName: projectName})
	} else {
		pq := query.NewProjectQuery()
		p, err = pq.GetMemberProject(userId, ownername, projectName)
	}

	if err != nil {
		return p, errors.New("validateProjectNameAndGetProject")
	}

	return p, nil
}

func validateTableIdAndGetTable (projectId, tableId int) (model.Table, error) {
	tr := repository.NewTableRepository()
	t, err := tr.GetOne(&model.Table{TableId: tableId})

	if err != nil || t.ProjectId != projectId {
		return t, errors.New("validateTableIdAndGetTable")
	}
	return t, nil
}

func validateColumnIdAndGetColumn (tableId, columnId int) (model.Column, error) {
	cr := repository.NewColumnRepository()
	c, err := cr.GetOne(&model.Column{ColumnId: columnId})
	
	if err != nil || c.TableId != tableId {
		return c, errors.New("validateColumnIdAndGetColumn")
	}
	return c, nil
}
package projects

import (
	"github.com/gin-gonic/gin"
	"task-management2/internal/controller/http/v1/projects"
)

func Router(g *gin.RouterGroup, projectsController *projects.Controller) {
	userG := g.Group("/projects")
	{
		// get-list
		userG.GET("/list", projectsController.GetProjectsWithStats)
		// get-detail
		userG.GET("/:id", projectsController.ProjectGetDetail)
		// create
		userG.POST("/create", projectsController.ProjectCreate)
		// update
		userG.PUT("/:id", projectsController.ProjectUpdate)
		// delete
		userG.DELETE("/:id", projectsController.ProjectDelete)

	}
}

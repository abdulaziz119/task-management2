package tasks

import (
	"github.com/gin-gonic/gin"
	"task-management2/internal/controller/http/v1/projects"
)

func Router(g *gin.RouterGroup, tasksController *projects.Controller) {
	userG := g.Group("/task")
	{
		// get-list
		userG.GET("/list", tasksController.ProjectGetList)
		// get-detail
		userG.GET("/:id", tasksController.ProjectGetDetail)
		// create
		userG.POST("/create", tasksController.ProjectCreate)
		// update
		userG.PUT("/:id", tasksController.ProjectUpdate)
		// delete
		userG.DELETE("/:id", tasksController.ProjectDelete)

	}
}

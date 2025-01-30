package tasks

import (
	"github.com/gin-gonic/gin"
	"task-management2/internal/controller/http/v1/tasks"
)

func Router(g *gin.RouterGroup, tasksController *tasks.Controller) {
	userG := g.Group("/task")
	{
		// get-list
		userG.GET("/list", tasksController.GetList)
		// get-detail
		userG.GET("/:id", tasksController.GetDetail)
		// create
		userG.POST("/create", tasksController.Create)
		// update
		userG.PUT("/:id", tasksController.Update)
		// delete
		userG.DELETE("/:id", tasksController.Delete)
	}
}

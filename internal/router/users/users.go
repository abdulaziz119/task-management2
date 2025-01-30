package users

import (
	"github.com/gin-gonic/gin"
	"task-management2/internal/controller/http/v1/users"
)

func Router(g *gin.RouterGroup, userController *users.Controller) {
	userG := g.Group("/user")
	{
		// get-list
		userG.GET("/list", userController.GetList)
		// get-detail
		userG.GET("/:id", userController.GetDetail)
		// create
		userG.POST("/create", userController.Create)
		// update
		userG.PUT("/:id", userController.Update)
		// delete
		userG.DELETE("/:id", userController.Delete)

	}
}

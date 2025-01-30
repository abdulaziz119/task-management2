package users

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	basic_controller "task-management2/internal/controller/http/v1/_basic_controller"
	"task-management2/internal/repository/postgres/users"
)

type Controller struct {
	useCase Repository
}

func NewController(useCase Repository) *Controller {
	return &Controller{useCase: useCase}
}

func (cl Controller) GetList(c *gin.Context) {
	var filter users.Filter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	list, count, err := cl.useCase.GetAll(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  list,
		"count": count,
	})
}

func (cl Controller) GetDetail(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "id must be a number!",
			"status":  false,
		})

		return
	}

	ctx := context.Background()

	detail, err := cl.useCase.GetById(ctx, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
		"data":    detail,
	})
}

func (cl Controller) Create(c *gin.Context) {
	var data users.Create
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := cl.useCase.Create(c.Request.Context(), data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (cl Controller) Update(c *gin.Context) {
	var data users.Update
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := cl.useCase.Update(c.Request.Context(), data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (cl Controller) Delete(c *gin.Context) {
	ctx, data, err := basic_controller.BasicDelete(c)
	if err != nil {
		return
	}

	err = cl.useCase.Delete(ctx, data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
	})
}

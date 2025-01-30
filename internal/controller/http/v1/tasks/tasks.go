package tasks

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	basic_controller "task-management2/internal/controller/http/v1/_basic_controller"
	"task-management2/internal/repository/postgres/tasks"
)

type Controller struct {
	useCase Repository
}

func NewController(useCase Repository) *Controller {
	return &Controller{useCase: useCase}
}

func (cl *Controller) GetList(c *gin.Context) {
	var filter tasks.Filter
	query := c.Request.URL.Query()

	defaultOffset := 0
	defaultLimit := 10
	filter.Offset = &defaultOffset
	filter.Limit = &defaultLimit

	projectIdQ := query["project_id"]
	if len(projectIdQ) > 0 {
		queryInt, err := strconv.Atoi(projectIdQ[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "project_id must be integer!",
				"status":  false,
			})
			return
		}
		filter.ProjectId = &queryInt
	}

	limitQ := query["limit"]
	if len(limitQ) > 0 {
		queryInt, err := strconv.Atoi(limitQ[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "limit must be number!",
				"status":  false,
			})
			return
		}
		filter.Limit = &queryInt
	}

	offsetQ := query["offset"]
	if len(offsetQ) > 0 {
		page, err := strconv.Atoi(offsetQ[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "offset must be number!",
				"status":  false,
			})
			return
		}
		offset := (page - 1) * *filter.Limit
		filter.Offset = &offset
	}

	list, taskStats, count, err := cl.useCase.GetAll(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       list,
		"count":      count,
		"task_stats": taskStats,
	})
}

func (cl *Controller) GetDetail(c *gin.Context) {
	var uri tasks.DetailUri

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	detail, err := cl.useCase.GetById(c.Request.Context(), uri.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": detail,
	})
}

func (cl *Controller) Create(c *gin.Context) {
	var request tasks.Create

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	detail, err := cl.useCase.Create(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": detail,
	})
}

func (cl *Controller) Update(c *gin.Context) {
	var uri tasks.DetailUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request tasks.Update
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uri.Id
	request.Id = &id

	detail, err := cl.useCase.Update(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": detail,
	})
}

func (cl *Controller) Delete(c *gin.Context) {
	ctx, data, err := basic_controller.BasicDelete(c)
	if err != nil {
		return
	}

	err = cl.useCase.Delete(ctx, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
	})
}

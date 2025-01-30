package export

import (
	"github.com/gin-gonic/gin"
	"task-management/internal/controller/http/v1/export"
)

func Router(g *gin.RouterGroup, exportController *export.Controller) {
	exportG := g.Group("/export")
	{
		exportG.GET("/excel", exportController.ExportToExcel)
	}
}

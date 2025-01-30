package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"

	export_controller "task-management2/internal/controller/http/v1/export"
	projects_controller "task-management2/internal/controller/http/v1/projects"
	tasks_controller "task-management2/internal/controller/http/v1/tasks"
	users_controller "task-management2/internal/controller/http/v1/users"
	"task-management2/internal/pkg/config"
	"task-management2/internal/pkg/repository/postgres"
	"task-management2/internal/repository/postgres/projects"
	"task-management2/internal/repository/postgres/tasks"
	"task-management2/internal/repository/postgres/users"
	"task-management2/internal/router/export"
	project_router "task-management2/internal/router/projects"
	task_router "task-management2/internal/router/tasks"
	user_router "task-management2/internal/router/users"
)

func main() {
	r := gin.Default()
	r.MaxMultipartMemory = 16 << 20

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	postgresDB := postgres.NewPostgres()

	// Repository
	userRepo := users.NewRepository(postgresDB)
	taskRepo := tasks.NewRepository(postgresDB)
	projectRepo := projects.NewRepository(postgresDB)

	// Controllers
	userController := users_controller.NewController(userRepo)
	taskController := tasks_controller.NewController(taskRepo)
	projectsController := projects_controller.NewController(projectRepo)
	exportController := export_controller.NewController(userRepo, taskRepo, projectRepo)

	api := r.Group("api")
	{
		v1 := api.Group("v1")

		v1.GET("/time", func(c *gin.Context) {
			now := time.Now()
			c.JSON(http.StatusOK, gin.H{
				"message": "ok!",
				"status":  true,
				"data": map[string]interface{}{
					"time":            now.Format("15:04"),
					"time_in_seconds": now.Hour()*3600 + now.Minute()*60 + now.Second(),
					"unix":            now.Unix(),
					"date":            now.Format("02.01.2006"),
					"week_day":        now.Weekday(),
					"full_date":       now.Format("02.01.2006 15:04:06"),
					"month":           now.Month(),
					"day":             now.Day(),
					"year":            now.Year(),
					"hour":            now.Hour(),
					"minute":          now.Minute(),
					"second":          now.Second(),
				},
			})
		})

		// Routers
		user_router.Router(v1, userController)
		task_router.Router(v1, taskController)
		project_router.Router(v1, projectsController)
		export.Router(v1, exportController)
	}

	log.Fatalln(r.Run(":" + config.GetConf().Port))
}

package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	export4 "task-management/internal/controller/http/v1/export"
	projects4 "task-management/internal/controller/http/v1/projects"
	tasks4 "task-management/internal/controller/http/v1/tasks"
	users4 "task-management/internal/controller/http/v1/users"
	"task-management/internal/pkg/config"
	"task-management/internal/pkg/repository/postgres"
	"task-management/internal/repository/postgres/projects"
	"task-management/internal/repository/postgres/tasks"
	"task-management/internal/repository/postgres/users"
	"task-management/internal/router/export"
	project_router "task-management/internal/router/projects"
	task_router "task-management/internal/router/tasks"
	user_router "task-management/internal/router/users"
	projects2 "task-management/internal/service/projects"
	tasks2 "task-management/internal/service/tasks"
	users2 "task-management/internal/service/users"
	projects3 "task-management/internal/usecase/projects"
	tasks3 "task-management/internal/usecase/tasks"
	users3 "task-management/internal/usecase/users"
	"time"
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

	//repository
	userRepo := users.NewRepository(postgresDB)
	taskRepo := tasks.NewRepository(postgresDB)
	projectRepo := projects.NewRepository(postgresDB)

	////service
	userService := users2.NewService(userRepo)
	tasksService := tasks2.NewService(taskRepo)
	projectsService := projects2.NewService(projectRepo)

	userUseCase := users3.NewUseCase(userService, userRepo)
	taskUseCaseService := tasks3.NewUseCase(tasksService)
	projectsUseCaseService := projects3.NewUseCase(projectsService)

	userController := users4.NewController(userUseCase)
	taskController := tasks4.NewController(taskUseCaseService)
	projectsController := projects4.NewController(projectsUseCaseService)
	exportController := export4.NewController(userUseCase, taskUseCaseService, projectsUseCaseService)

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
		// @router
		{
			user_router.Router(v1, userController)
			task_router.Router(v1, taskController)
			project_router.Router(v1, projectsController)
			export.Router(v1, exportController)
		}
	}

	log.Fatalln(r.Run(":" + config.GetConf().Port))
}

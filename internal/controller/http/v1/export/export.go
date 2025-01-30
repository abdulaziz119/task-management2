package export

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"task-management/internal/service/projects"
	"task-management/internal/service/tasks"
	"task-management/internal/service/users"
	usecase_projects "task-management/internal/usecase/projects"
	usecase_tasks "task-management/internal/usecase/tasks"
	usecase_users "task-management/internal/usecase/users"
)

type Controller struct {
	userUseCase    *usecase_users.UseCase
	taskUseCase    *usecase_tasks.TasksUseCase
	projectUseCase *usecase_projects.UseCase
}

func NewController(userUseCase *usecase_users.UseCase, taskUseCase *usecase_tasks.TasksUseCase, projectUseCase *usecase_projects.UseCase) *Controller {
	return &Controller{
		userUseCase:    userUseCase,
		taskUseCase:    taskUseCase,
		projectUseCase: projectUseCase,
	}
}

func (h *Controller) ExportToExcel(c *gin.Context) {
	ctx := c.Request.Context()

	userList, _, err := h.userUseCase.GetAll(ctx, users.Filter{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Error getting users: %v", err),
		})
		return
	}

	userMap := make(map[int64]string)
	for _, u := range userList {
		if u.Id != nil {
			fullName := ""
			if u.FullName != nil {
				fullName = *u.FullName
			}
			userMap[*u.Id] = fullName
		}
	}

	projectList, _, err := h.projectUseCase.ProjectGetList(ctx, projects.Filter{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Error getting projects: %v", err),
		})
		return
	}

	projectMap := make(map[int]string)
	for _, p := range projectList {
		projectMap[p.Id] = p.Name
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	userSheet := "Users"
	f.NewSheet(userSheet)
	f.SetCellValue(userSheet, "A1", "ID")
	f.SetCellValue(userSheet, "B1", "Full Name")
	f.SetCellValue(userSheet, "C1", "Email")
	f.SetCellValue(userSheet, "D1", "Role")
	f.SetCellValue(userSheet, "E1", "Pending Tasks")
	f.SetCellValue(userSheet, "F1", "In Progress Tasks")
	f.SetCellValue(userSheet, "G1", "Completed Tasks")
	f.SetCellValue(userSheet, "H1", "Total Tasks")

	for i, u := range userList {
		row := i + 2
		f.SetCellValue(userSheet, fmt.Sprintf("A%d", row), *u.Id)
		f.SetCellValue(userSheet, fmt.Sprintf("B%d", row), *u.FullName)
		f.SetCellValue(userSheet, fmt.Sprintf("C%d", row), *u.Email)
		f.SetCellValue(userSheet, fmt.Sprintf("D%d", row), *u.Role)
		f.SetCellValue(userSheet, fmt.Sprintf("E%d", row), *u.PendingTasks)
		f.SetCellValue(userSheet, fmt.Sprintf("F%d", row), *u.InProgressTasks)
		f.SetCellValue(userSheet, fmt.Sprintf("G%d", row), *u.CompletedTasks)
		total := *u.PendingTasks + *u.InProgressTasks + *u.CompletedTasks
		f.SetCellValue(userSheet, fmt.Sprintf("H%d", row), total)
	}

	tasksList, _, _, err := h.taskUseCase.TaskGetList(ctx, tasks.Filter{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Error getting tasks: %v", err),
		})
		return
	}

	taskSheet := "Tasks"
	f.NewSheet(taskSheet)
	f.SetCellValue(taskSheet, "A1", "ID")
	f.SetCellValue(taskSheet, "B1", "Name")
	f.SetCellValue(taskSheet, "C1", "Description")
	f.SetCellValue(taskSheet, "D1", "Project")
	f.SetCellValue(taskSheet, "E1", "Status")
	f.SetCellValue(taskSheet, "F1", "Priority")
	f.SetCellValue(taskSheet, "G1", "Due Date")
	f.SetCellValue(taskSheet, "H1", "Assigned To")

	redStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FF9999"}, Pattern: 1},
	})
	yellowStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FFEB9C"}, Pattern: 1},
	})
	greenStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#C6EFCE"}, Pattern: 1},
	})

	for i, task := range tasksList {
		row := i + 2
		f.SetCellValue(taskSheet, fmt.Sprintf("A%d", row), task.Id)
		f.SetCellValue(taskSheet, fmt.Sprintf("B%d", row), task.Name)
		f.SetCellValue(taskSheet, fmt.Sprintf("C%d", row), task.Description)

		if projectName, ok := projectMap[task.ProjectId]; ok {
			f.SetCellValue(taskSheet, fmt.Sprintf("D%d", row), projectName)
		} else {
			f.SetCellValue(taskSheet, fmt.Sprintf("D%d", row), task.ProjectId)
		}

		f.SetCellValue(taskSheet, fmt.Sprintf("E%d", row), task.Status)
		f.SetCellValue(taskSheet, fmt.Sprintf("F%d", row), task.Priority)
		f.SetCellValue(taskSheet, fmt.Sprintf("G%d", row), task.DueDate)

		assignedTo := int64(task.AssignedTo)
		if assignedTo > 0 {
			if userName, exists := userMap[assignedTo]; exists && userName != "" {
				f.SetCellValue(taskSheet, fmt.Sprintf("H%d", row), userName)
			} else {
				f.SetCellValue(taskSheet, fmt.Sprintf("H%d", row), "Not assigned")
			}
		} else {
			f.SetCellValue(taskSheet, fmt.Sprintf("H%d", row), "Not assigned")
		}

		statusCell := fmt.Sprintf("E%d", row)
		switch task.Status {
		case "pending":
			f.SetCellStyle(taskSheet, statusCell, statusCell, redStyle)
		case "in_progress":
			f.SetCellStyle(taskSheet, statusCell, statusCell, yellowStyle)
		case "completed":
			f.SetCellStyle(taskSheet, statusCell, statusCell, greenStyle)
		}
	}

	projectSheet := "Projects"
	f.NewSheet(projectSheet)
	f.SetCellValue(projectSheet, "A1", "ID")
	f.SetCellValue(projectSheet, "B1", "Name")
	f.SetCellValue(projectSheet, "C1", "Description")
	f.SetCellValue(projectSheet, "D1", "Owner ID")
	f.SetCellValue(projectSheet, "E1", "Total Tasks")
	f.SetCellValue(projectSheet, "F1", "Progress")

	for i, p := range projectList {
		row := i + 2
		f.SetCellValue(projectSheet, fmt.Sprintf("A%d", row), p.Id)
		f.SetCellValue(projectSheet, fmt.Sprintf("B%d", row), p.Name)
		f.SetCellValue(projectSheet, fmt.Sprintf("C%d", row), p.Description)
		f.SetCellValue(projectSheet, fmt.Sprintf("D%d", row), p.OwnerId)
		f.SetCellValue(projectSheet, fmt.Sprintf("E%d", row), p.TotalTasks)
		f.SetCellValue(projectSheet, fmt.Sprintf("F%d", row), fmt.Sprintf("%.2f%%", p.Progress))
	}

	f.SetActiveSheet(0)

	filename := fmt.Sprintf("task_management_export_%s.xlsx", time.Now().Format("2006-01-02_15-04-05"))

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Error writing file: %v", err),
		})
		return
	}
}

func (h *Controller) ExportProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := h.projectUseCase.ProjectGetDetail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

func getValueOrEmpty(v interface{}) interface{} {
	if v == nil {
		return ""
	}
	return v
}

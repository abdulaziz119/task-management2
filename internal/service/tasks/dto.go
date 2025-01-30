package tasks

type Filter struct {
	Limit     *int
	Offset    *int
	ProjectId *int
}

type Create struct {
	ProjectId   *int    `json:"project_id" bun:"project_id"`
	Name        *string `json:"name" bun:"name"`
	Description *string `json:"description" bun:"description"`
	AssignedTo  *int    `json:"assigned_to" bun:"assigned_to"`
	Status      *string `json:"status" validate:"required,oneof=pending in_progress completed"`
	Priority    *string `json:"priority" validate:"required,oneof=low medium high"`
	DueDate     *string `json:"due_date" bun:"due_date"`
}

type Update struct {
	Id          *int    `json:"id" form:"id"`
	ProjectId   *int    `json:"project_id" bun:"project_id"`
	Name        *string `json:"name" bun:"name"`
	Description *string `json:"description" bun:"description"`
	AssignedTo  *int    `json:"assigned_to" bun:"assigned_to"`
	Status      *string `json:"status" validate:"required,oneof=pending in_progress completed"`
	Priority    *string `json:"priority" validate:"required,oneof=low medium high"`
	DueDate     *string `json:"due_date" bun:"due_date"`
}

type TaskStats struct {
	TotalTasks      int     `json:"total_tasks"`
	CompletedTasks  int     `json:"completed_tasks"`
	InProgressTasks int     `json:"in_progress_tasks"`
	PendingTasks    int     `json:"pending_tasks"`
	Progress        float64 `json:"progress"`
}

type List struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ProjectId   int    `json:"project_id"`
	AssignedTo  int    `json:"assigned_to"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
}

type DetailUri struct {
	Id int `uri:"id" binding:"required"`
}

type Detail struct {
	Id          int     `json:"id"`
	ProjectId   *int    `json:"project_id" bun:"project_id"`
	Name        *string `json:"name" bun:"name"`
	Description *string `json:"description" bun:"description"`
	AssignedTo  *int    `json:"assigned_to" bun:"assigned_to"`
	Status      *string `json:"status" bun:"status"`
	Priority    *string `json:"priority" bun:"priority"`
	DueDate     *string `json:"due_date" bun:"due_date"`
}

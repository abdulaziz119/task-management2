package projects

type Filter struct {
	Limit   *int
	Offset  *int
	OwnerId *int
}

type Create struct {
	Name        *string `json:"name" bun:"name"`
	Description *string `json:"description" bun:"description"`
	Owner_id    *int    `json:"owner_id" bun:"owner_id"`
}

type Update struct {
	Id          *int    `json:"id" form:"id"`
	Name        *string `json:"name" bun:"name"`
	Description *string `json:"description" bun:"description"`
	Owner_id    *int    `json:"owner_id" bun:"owner_id"`
}

type TaskStats struct {
	TotalTasks int     `json:"total_tasks"`
	Progress   float64 `json:"progress"`
}

type List struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	OwnerId     int     `json:"owner_id"`
	TotalTasks  int     `json:"total_tasks"`
	Progress    float64 `json:"progress"`
}

type Detail struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Owner_id    int       `json:"owner_id"`
	TaskStats   TaskStats `json:"task_stats"`
}

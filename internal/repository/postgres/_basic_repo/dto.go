package basic_repo

type Delete struct {
	Id *int `json:"id" form:"id" bun:"id"`
}

package basic_service

type Delete struct {
	Id *int `json:"id" form:"id" bun:"id"`
}

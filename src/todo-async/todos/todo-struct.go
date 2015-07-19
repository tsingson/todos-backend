package todos

type Todo struct {
	Id          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description" binding:"required"`
}

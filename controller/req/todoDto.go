package req

type CreateTodoDto struct {
	Title    string `json:"title" example:"제목"`
	Content  string `json:"content" example:"내용"`
	OrderNum int    `json:"order_num" example:"2"`
}

type UpdateTodoDto struct {
	Id        int    `json:"id" example:"0"`
	Title     string `json:"title" example:"제목"`
	Content   string `json:"content" example:"내용"`
	OrderNum  int    `json:"order_num" example:"1"`
	IsDeleted bool   `json:"is_deleted" example:"false"`
}

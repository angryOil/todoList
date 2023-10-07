package req

import (
	"todoList/domain"
)

// id
// title
// content
// createdAt
// orderNum
// isDeleted

type CreateTodoDto struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	OrderNum int    `json:"order_num"`
}

func (d CreateTodoDto) ToDomain() domain.Todo {
	return domain.Todo{
		Title:    d.Title,
		Content:  d.Content,
		OrderNum: d.OrderNum,
	}
}

type UpdateTodoDto struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	OrderNum  int    `json:"order_num"`
	IsDeleted bool   `json:"is_deleted"`
}

func (d UpdateTodoDto) ToDomain() domain.Todo {
	return domain.Todo{
		Id:        d.Id,
		Title:     d.Title,
		Content:   d.Content,
		OrderNum:  d.OrderNum,
		IsDeleted: d.IsDeleted,
	}
}

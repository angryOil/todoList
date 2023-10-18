package req

import (
	"time"
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

func (d CreateTodoDto) ToDomain(userId int) domain.Todo {
	return domain.Todo{
		UserId:   userId,
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

func (d UpdateTodoDto) ToDomain(userId int) domain.Todo {
	return domain.Todo{
		Id:            d.Id,
		UserId:        userId,
		Title:         d.Title,
		Content:       d.Content,
		OrderNum:      d.OrderNum,
		LastUpdatedAt: time.Now(),
		IsDeleted:     d.IsDeleted,
	}
}

package req

import (
	"time"
	"todoList/domain"
)

type CreateTodoDto struct {
	Title    string `json:"title" example:"제목"`
	Content  string `json:"content" example:"내용"`
	OrderNum int    `json:"order_num" example:"2"`
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
	Id        int    `json:"id" example:"0"`
	Title     string `json:"title" example:"제목"`
	Content   string `json:"content" example:"내용"`
	OrderNum  int    `json:"order_num" example:"1"`
	IsDeleted bool   `json:"is_deleted" example:"false"`
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

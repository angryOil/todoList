package model

import (
	"github.com/uptrace/bun"
	"time"
	"todoList/domain"
	"todoList/repository/req"
)

type Todo struct {
	bun.BaseModel `bun:"table:todo,alias:t"`

	Id            int       `bun:"id,pk,autoincrement"`
	UserId        int       `bun:"user_id,bigint,notnull"`
	Title         string    `bun:"title,notnull"`
	Content       string    `bun:"content"`
	OrderNum      int       `bun:"order_num,notnull"`
	IsDeleted     bool      `bun:"is_deleted,notnull"`
	CreatedAt     time.Time `bun:"created_at,notnull"`
	LastUpdatedAt time.Time `bun:"last_updated_at,notnull"`
}

func ToCreateModel(c req.CreateTodo) Todo {
	return Todo{
		UserId:    c.UserId,
		Title:     c.Title,
		Content:   c.Content,
		OrderNum:  c.OrderNum,
		IsDeleted: c.IsDeleted,
		CreatedAt: c.CreatedAt,
	}
}

func ToSaveModel(s req.Save) Todo {
	return Todo{
		Id:            s.Id,
		UserId:        s.UserId,
		Title:         s.Title,
		Content:       s.Content,
		OrderNum:      s.OrderNum,
		IsDeleted:     s.IsDeleted,
		CreatedAt:     s.CreatedAt,
		LastUpdatedAt: s.LastUpdatedAt,
	}
}

func ToDomainList(list []Todo) []domain.Todo {
	result := make([]domain.Todo, len(list))

	for i, t := range list {
		result[i] = domain.NewTodoBuilder().
			Id(t.Id).
			UserId(t.UserId).
			Title(t.Title).
			Content(t.Content).
			OrderNum(t.OrderNum).
			IsDeleted(t.IsDeleted).
			CreatedAt(t.CreatedAt).
			LastUpdatedAt(t.LastUpdatedAt).
			Build()
	}
	return result
}

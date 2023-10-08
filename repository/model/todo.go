package model

import (
	"github.com/uptrace/bun"
	"time"
	"todoList/domain"
)

type Todo struct {
	bun.BaseModel `bun:"table:todo,alias:t"`

	Id            int       `bun:"id,pk,autoincrement"`
	Title         string    `bun:"title,notnull"`
	OrderNum      int       `bun:"order_num,notnull"`
	IsDeleted     bool      `bun:"is_deleted,notnull"`
	CreatedAt     time.Time `bun:"created_at,notnull"`
	LastUpdatedAt time.Time `bun:"last_updated_at,notnull"`
}

func (t Todo) ToDomain() domain.Todo {
	return domain.Todo{
		Id:            t.Id,
		Title:         t.Title,
		OrderNum:      t.OrderNum,
		IsDeleted:     t.IsDeleted,
		CreatedAt:     t.CreatedAt,
		LastUpdatedAt: t.LastUpdatedAt,
	}
}

type TodoDetail struct {
	bun.BaseModel `bun:"table:todo,alias:t"`

	Id            int       `bun:"id,pk,autoincrement"`
	Title         string    `bun:"title,notnull"`
	Content       string    `bun:"content,notnull"`
	OrderNum      int       `bun:"order_num,notnull"`
	IsDeleted     bool      `bun:"is_deleted,notnull"`
	CreatedAt     time.Time `bun:"created_at,notnull"`
	LastUpdatedAt time.Time `bun:"last_updated_at,notnull"`
}

func (t TodoDetail) ToDomain() domain.Todo {
	return domain.Todo{
		Id:            t.Id,
		Title:         t.Title,
		Content:       t.Content,
		OrderNum:      t.OrderNum,
		IsDeleted:     t.IsDeleted,
		CreatedAt:     t.CreatedAt,
		LastUpdatedAt: t.LastUpdatedAt,
	}
}

func ToDomainList(list []Todo) []domain.Todo {
	result := make([]domain.Todo, len(list))

	for i, todo := range list {
		result[i] = todo.ToDomain()
	}
	return result
}

func ToDomainDetailList(list []TodoDetail) []domain.Todo {
	result := make([]domain.Todo, len(list))

	for i, todo := range list {
		result[i] = todo.ToDomain()
	}
	return result
}

func ToDetailModel(dt domain.Todo) TodoDetail {
	return TodoDetail{
		Id:            dt.Id,
		Title:         dt.Title,
		Content:       dt.Content,
		OrderNum:      dt.OrderNum,
		IsDeleted:     dt.IsDeleted,
		CreatedAt:     dt.CreatedAt,
		LastUpdatedAt: dt.LastUpdatedAt,
	}
}

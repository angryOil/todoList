package repository

import (
	"context"
	"todoList/domain"
)

type ITodoRepository interface {
	Create(todo domain.Todo) error
	Save(todo domain.Todo, saveFunc func(todo2 domain.Todo) error) error
	GetDetail(ctx context.Context, id int) ([]domain.Todo, error)
	GetList(ctx context.Context) ([]domain.Todo, error)
}

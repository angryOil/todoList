package repository

import (
	"context"
	"todoList/domain"
	"todoList/page"
)

type ITodoRepository interface {
	Create(ctx context.Context, todo domain.Todo) error
	Save(ctx context.Context, todo domain.Todo, saveFunc func(todo2 domain.Todo) error) error
	Delete(ctx context.Context, userId, id int) error
	GetDetail(ctx context.Context, userId, id int) ([]domain.Todo, error)
	GetList(ctx context.Context, userId int, page page.ReqPage) ([]domain.Todo, int, error)
}

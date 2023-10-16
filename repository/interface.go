package repository

import (
	"context"
	"todoList/domain"
	"todoList/page"
)

type ITodoRepository interface {
	Create(ctx context.Context, todo domain.Todo) error
	Save(ctx context.Context, userId, id int,
		getValidFunc func([]domain.Todo) (domain.Todo, error),
		mergeTodo func(todo domain.Todo) domain.Todo,
		saveValidFunc func(domain.Todo) error) error
	Delete(ctx context.Context, userId, id int) error
	GetDetail(ctx context.Context, userId, id int) ([]domain.Todo, error)
	GetList(ctx context.Context, userId int, page page.ReqPage) ([]domain.Todo, int, error)
}

package repository

import (
	"context"
	"todoList/domain"
	"todoList/domain/vo"
	"todoList/page"
	"todoList/repository/req"
)

type ITodoRepository interface {
	Create(ctx context.Context, todo req.CreateTodo) error
	Save(ctx context.Context, userId, id int,
		getValidFunc func([]domain.Todo) (domain.Todo, error),
		mergeTodo func(todo domain.Todo) (vo.Save, error)) error
	Delete(ctx context.Context, userId, id int) error
	GetDetail(ctx context.Context, userId, id int) ([]domain.Todo, error)
	GetList(ctx context.Context, userId int, page page.ReqPage) ([]domain.Todo, int, error)
}

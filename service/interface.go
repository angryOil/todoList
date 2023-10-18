package service

import (
	"context"
	"todoList/domain"
	"todoList/page"
)

type ITodoService interface {
	CreateTodo(ctx context.Context, todo domain.Todo) error
	UpdateTodo(ctx context.Context, todo domain.Todo) error
	DeleteTodo(ctx context.Context, userId, id int) error
	GetTodos(ctx context.Context, user int, page page.ReqPage) ([]domain.Todo, int, error)
	GetDetail(ctx context.Context, userId, id int) (domain.Todo, error)
}

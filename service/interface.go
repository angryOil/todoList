package service

import (
	"context"
	"todoList/domain"
)

type ITodoService interface {
	CreateTodo(ctx context.Context, todo domain.Todo) error
	UpdateTodo(ctx context.Context, todo domain.Todo) error
	DeleteTodo(ctx context.Context, id int) error
	GetTodos(ctx context.Context) ([]domain.Todo, error)
	GetDetail(ctx context.Context, id int) (domain.Todo, error)
}

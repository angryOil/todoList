package service

import (
	"context"
	"todoList/page"
	"todoList/service/req"
	"todoList/service/res"
)

type ITodoService interface {
	CreateTodo(ctx context.Context, c req.CreateTodo) error
	UpdateTodo(ctx context.Context, s req.Save) error
	DeleteTodo(ctx context.Context, userId, id int) error
	GetTodos(ctx context.Context, user int, page page.ReqPage) ([]res.GetList, int, error)
	GetDetail(ctx context.Context, userId, id int) (res.GetDetail, error)
}

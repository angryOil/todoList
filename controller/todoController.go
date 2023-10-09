package controller

import (
	"context"
	"todoList/controller/req"
	"todoList/controller/res"
	"todoList/service"
)

// dto/domain 을 둘다 사용 res/req 는 dto , service 호출/응답 은 domain 으로 통신

type TodoController struct {
	service service.ITodoService
}

func NewController(serv service.ITodoService) TodoController {
	return TodoController{
		service: serv,
	}
}

func (c TodoController) CreateTodo(ctx context.Context, dto req.CreateTodoDto) error {
	err := c.service.CreateTodo(ctx, dto.ToDomain())
	if err != nil {
		return err
	}

	return nil
}

func (c TodoController) UpdateTodo(ctx context.Context, dto req.UpdateTodoDto) error {
	err := c.service.UpdateTodo(ctx, dto.ToDomain())
	return err
}

// todo transaction 익힌후 테스트

func (c TodoController) GetTodos(ctx context.Context) ([]res.ListDto, error) {
	todoDomains, err := c.service.GetTodos(ctx)
	if err != nil {
		return []res.ListDto{}, err
	}

	return res.ToListDtoList(todoDomains), nil
}

func (c TodoController) GetDetail(ctx context.Context, id int) (res.DetailDto, error) {
	todo, err := c.service.GetDetail(ctx, id)
	if err != nil {
		return res.DetailDto{}, err
	}
	return res.ToDetailDto(todo), nil
}

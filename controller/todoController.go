package controller

import (
	"context"
	"errors"
	"todoList/controller/req"
	"todoList/controller/res"
	"todoList/page"
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
	userId, ok := ctx.Value("userId").(int)
	if !ok {
		return errors.New("user id is not valid")
	}
	err := c.service.CreateTodo(ctx, dto.ToDomain(userId))
	if err != nil {
		return err
	}

	return nil
}

func (c TodoController) UpdateTodo(ctx context.Context, dto req.UpdateTodoDto) error {
	userId, ok := ctx.Value("userId").(int)
	if !ok {
		return errors.New("user id is not valid")
	}
	err := c.service.UpdateTodo(ctx, dto.ToDomain(userId))
	return err
}

func (c TodoController) DeleteTodo(ctx context.Context, id int) error {
	userId, ok := ctx.Value("userId").(int)
	if !ok {
		return errors.New("user id is not valid")
	}
	err := c.service.DeleteTodo(ctx, userId, id)
	return err
}

func (c TodoController) GetTodos(ctx context.Context, page page.ReqPage) ([]res.ListDto, int, error) {
	userId, ok := ctx.Value("userId").(int)
	if !ok {
		return []res.ListDto{}, 0, errors.New("user id is not valid")
	}
	todoDomains, count, err := c.service.GetTodos(ctx, userId, page)
	if err != nil {
		return []res.ListDto{}, 0, err
	}

	return res.ToListDtoList(todoDomains), count, nil
}

func (c TodoController) GetDetail(ctx context.Context, id int) (res.DetailDto, error) {
	userId, ok := ctx.Value("userId").(int)
	if !ok {
		return res.DetailDto{}, errors.New("user id is not valid")
	}
	todo, err := c.service.GetDetail(ctx, userId, id)
	if err != nil {
		return res.DetailDto{}, err
	}
	return res.ToDetailDto(todo), nil
}

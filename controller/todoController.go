package controller

import (
	"context"
	"errors"
	"time"
	"todoList/controller/req"
	"todoList/controller/res"
	"todoList/page"
	"todoList/service"
	req2 "todoList/service/req"
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
	err := c.service.CreateTodo(ctx, req2.CreateTodo{
		UserId:   userId,
		Title:    dto.Title,
		Content:  dto.Content,
		OrderNum: dto.OrderNum,
	})

	return err
}

func (c TodoController) UpdateTodo(ctx context.Context, dto req.UpdateTodoDto) error {
	userId, ok := ctx.Value("userId").(int)
	if !ok {
		return errors.New("user id is not valid")
	}
	err := c.service.UpdateTodo(ctx, req2.Save{
		Id:        dto.Id,
		UserId:    userId,
		Title:     dto.Title,
		Content:   dto.Content,
		OrderNum:  dto.OrderNum,
		IsDeleted: dto.IsDeleted,
	})
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
	listArr, count, err := c.service.GetTodos(ctx, userId, page)
	if err != nil {
		return []res.ListDto{}, 0, err
	}

	result := make([]res.ListDto, len(listArr))
	for i, l := range listArr {
		result[i] = res.ListDto{
			Id:        l.Id,
			Title:     l.Title,
			OrderNum:  l.OrderNum,
			CreatedAt: convertTimeToString(l.CreatedAt),
			IsDeleted: l.IsDeleted,
		}
	}
	return result, count, nil
}

func (c TodoController) GetDetail(ctx context.Context, id int) (res.DetailDto, error) {
	userId, ok := ctx.Value("userId").(int)
	if !ok {
		return res.DetailDto{}, errors.New("user id is not valid")
	}
	t, err := c.service.GetDetail(ctx, userId, id)
	if err != nil {
		return res.DetailDto{}, err
	}
	return res.DetailDto{
		Id:        t.Id,
		Title:     t.Title,
		UserId:    t.UserId,
		Content:   t.Content,
		CreatedAt: convertTimeToString(t.CreatedAt),
		OrderNum:  t.OrderNum,
		IsDeleted: t.IsDeleted,
	}, nil
}

var koreaZone, _ = time.LoadLocation("Asia/Seoul")

func convertTimeToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	t = t.In(koreaZone)
	return t.Format(time.RFC3339)
}

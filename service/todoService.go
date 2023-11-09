package service

import (
	"context"
	"errors"
	"time"
	"todoList/domain"
	"todoList/domain/vo"
	"todoList/page"
	"todoList/repository"
	req2 "todoList/repository/req"
	"todoList/service/req"
	"todoList/service/res"
)

// service domain => 정책/벨리데이션 => domain
// service 는 domain 만 사용

type TodoService struct {
	repo repository.ITodoRepository
}

func NewService(repo repository.ITodoRepository) TodoService {
	return TodoService{repo: repo}
}

const (
	NoRows        = "no rows"
	InvalidUserID = "invalid user id"
	InvalidId     = "invalid id"
)

func (s TodoService) CreateTodo(ctx context.Context, c req.CreateTodo) error {
	userId, orderNum := c.UserId, c.OrderNum
	title, content := c.Title, c.Content
	createdAt := time.Now()
	isDeleted := false

	err := domain.NewTodoBuilder().
		UserId(userId).
		OrderNum(orderNum).
		Title(title).
		Content(content).
		CreatedAt(createdAt).
		IsDeleted(isDeleted).
		Build().ValidCreate()
	if err != nil {
		return err
	}

	err = s.repo.Create(ctx, req2.CreateTodo{
		UserId:    userId,
		Title:     title,
		Content:   content,
		OrderNum:  orderNum,
		IsDeleted: isDeleted,
		CreatedAt: createdAt,
	})
	return err
}

func (s TodoService) DeleteTodo(ctx context.Context, userId, id int) error {
	err := s.repo.Delete(ctx, userId, id)
	return err
}

func (s TodoService) UpdateTodo(ctx context.Context, t req.Save) error {
	id, userId, orderNum := t.Id, t.UserId, t.OrderNum
	title, content := t.Title, t.Content
	isDeleted := t.IsDeleted

	err := s.repo.Save(ctx,
		userId, id,
		func(todos []domain.Todo) (domain.Todo, error) {
			if len(todos) == 0 {
				return domain.NewTodoBuilder().Build(), errors.New(NoRows)
			}
			return todos[0], nil
		},
		func(t domain.Todo) (vo.Save, error) {
			u := t.Update(title, content, orderNum, isDeleted)
			err := u.ValidUpdate()
			if err != nil {
				return vo.Save{}, err
			}
			return u.ToSave(), nil
		},
	)

	return err
}

func (s TodoService) GetTodos(ctx context.Context, userId int, page page.ReqPage) ([]res.GetList, int, error) {
	if userId == 0 {
		return []res.GetList{}, 0, errors.New(InvalidUserID)
	}
	todos, totalCount, err := s.repo.GetList(ctx, userId, page)
	result := make([]res.GetList, len(todos))
	for i, t := range todos {
		v := t.ToInfo()
		result[i] = res.GetList{
			Id:        v.Id,
			UserId:    v.UserId,
			Title:     v.Title,
			CreatedAt: v.CreatedAt,
			OrderNum:  v.OrderNum,
			IsDeleted: v.IsDeleted,
		}
	}
	return result, totalCount, err
}

func (s TodoService) GetDetail(ctx context.Context, userId, id int) (res.GetDetail, error) {
	if id == 0 {
		return res.GetDetail{}, errors.New(InvalidId)
	}
	detail, err := s.repo.GetDetail(ctx, userId, id)
	if err != nil {
		return res.GetDetail{}, err
	}
	if len(detail) == 0 {
		return res.GetDetail{}, nil
	}
	v := detail[0].ToDetail()
	return res.GetDetail{
		Id:        v.Id,
		Title:     v.Title,
		UserId:    v.UserId,
		Content:   v.Content,
		CreatedAt: v.CreatedAt,
		OrderNum:  v.OrderNum,
		IsDeleted: v.IsDeleted,
	}, nil
}

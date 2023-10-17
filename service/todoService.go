package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"todoList/domain"
	"todoList/page"
	"todoList/repository"
)

// service domain => 정책/벨리데이션 => domain
// service 는 domain 만 사용

type TodoService struct {
	repo repository.ITodoRepository
}

func NewService(repo repository.ITodoRepository) TodoService {
	return TodoService{repo: repo}
}
func (s TodoService) CreateTodo(ctx context.Context, todo domain.Todo) error {
	createdTodo, err := domain.CreatedTodo(todo.UserId, todo.Title, todo.Content, todo.OrderNum)
	if err != nil {
		return err
	}

	err = s.repo.Create(ctx, createdTodo)
	return err
}

func (s TodoService) DeleteTodo(ctx context.Context, userId, id int) error {
	err := s.repo.Delete(ctx, userId, id)
	return err
}

func (s TodoService) UpdateTodo(ctx context.Context, todo domain.Todo) error {
	err := s.repo.Save(ctx,
		todo.UserId, todo.Id,
		func(todos []domain.Todo) (domain.Todo, error) {
			if len(todos) == 0 {
				return domain.Todo{}, errors.New("no row error")
			}
			if todos[0].UserId != todo.UserId {
				return domain.Todo{}, errors.New("it`s not yours error who are u?")
			}
			return todos[0], nil
		},
		func(t domain.Todo) domain.Todo {
			t.Title = todo.Title
			t.Content = todo.Content
			t.IsDeleted = todo.IsDeleted
			t.OrderNum = todo.OrderNum
			t.LastUpdatedAt = time.Now()
			return t
		},
		updateValidFunc,
	)

	return err
}

func updateValidFunc(t domain.Todo) error {
	if t.Id == 0 {
		return errors.New("todoId is zero")
	}
	if t.UserId == 0 {
		return errors.New("userId is zero")
	}
	return nil
}

func (s TodoService) GetTodos(ctx context.Context, userId int, page page.ReqPage) ([]domain.Todo, int, error) {
	todos, totalCount, err := s.repo.GetList(ctx, userId, page)
	return todos, totalCount, err
}

func (s TodoService) GetDetail(ctx context.Context, userId, id int) (domain.Todo, error) {
	if id == 0 {
		return domain.Todo{}, errors.New("id is no value")
	}
	detail, err := s.repo.GetDetail(ctx, userId, id)
	if err != nil {
		return domain.Todo{}, err
	}
	if len(detail) == 0 {
		return domain.Todo{}, errors.New(fmt.Sprintf("no rows error: %d", id))
	}
	return detail[0], nil
}

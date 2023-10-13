package service

import (
	"context"
	"errors"
	"fmt"
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
	createdTodo, err1 := domain.CreatedTodo(todo.UserId, todo.Title, todo.Content, todo.OrderNum)
	if err1 != nil {
		return err1
	}

	err2 := s.repo.Create(ctx, createdTodo)
	return err2
}

func (s TodoService) DeleteTodo(ctx context.Context, userId, id int) error {
	err := s.repo.Delete(ctx, userId, id)
	return err
}

func (s TodoService) UpdateTodo(ctx context.Context, todo domain.Todo) error {
	todos, err := s.repo.GetDetail(ctx, todo.UserId, todo.Id)
	if err != nil {
		return err
	}
	if len(todos) == 0 {
		return errors.New("empty rows")
	}

	// select => update
	err = s.repo.Save(ctx, todo, func(todoDomain domain.Todo) error {
		if todo.Id == 0 {
			return errors.New("id is no value")
		}
		_, err := domain.CreatedTodo(todo.UserId, todo.Title, todo.Content, todo.OrderNum)
		if err != nil {
			return err
		}
		return nil
	})

	return err
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

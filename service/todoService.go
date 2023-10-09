package service

import (
	"context"
	"errors"
	"fmt"
	"todoList/domain"
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
	createdTodo, err1 := domain.CreatedTodo(todo.Title, todo.Content, todo.OrderNum)
	if err1 != nil {
		return err1
	}

	err2 := s.repo.Create(ctx, createdTodo)
	return err2
}

func (s TodoService) DeleteTodo(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	return err
}

func (s TodoService) UpdateTodo(ctx context.Context, todo domain.Todo) error {
	// select => update
	err := s.repo.Save(ctx, todo, func(todoDomain domain.Todo) error {
		if todo.Id == 0 {
			return errors.New("id is no value")
		}
		_, err := domain.CreatedTodo(todo.Title, todo.Content, todo.OrderNum)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

// todo transaction 익힌후 테스트

func (s TodoService) GetTodos(ctx context.Context) ([]domain.Todo, error) {
	todos, err := s.repo.GetList(ctx)
	return todos, err
}

func (s TodoService) GetDetail(ctx context.Context, id int) (domain.Todo, error) {
	if id == 0 {
		return domain.Todo{}, errors.New("id is no value")
	}
	detail, err := s.repo.GetDetail(ctx, id)
	if err != nil {
		return domain.Todo{}, err
	}
	if len(detail) == 0 {
		return domain.Todo{}, errors.New(fmt.Sprintf("no rows error: %d", id))
	}
	return detail[0], nil
}

package service

import (
	"context"
	"errors"
	"fmt"
	"todoList/domain"
	"todoList/repository"
)

// service domain => 정책/벨리데이션 => domain

type TodoService struct {
	repo repository.ITodoRepository
}

func (s TodoService) CreateTodo(ctx context.Context, todo domain.Todo) error {
	createdTodo, err1 := domain.CreatedTodo(todo.Title, todo.Content, todo.OrderNum)
	if err1 != nil {
		return err1
	}

	err2 := s.repo.Create(createdTodo)
	return err2
}

func (s TodoService) UpdateTodo(ctx context.Context, todo domain.Todo) error {
	// select => update
	err := s.repo.Save(todo, func(todoDomain domain.Todo) error {
		_, err := domain.CreatedTodo(todo.Title, todo.Content, todo.OrderNum)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s TodoService) GetTodos(ctx context.Context) ([]domain.Todo, error) {
	todos, err := s.repo.GetList(ctx)
	return todos, err
}

func (s TodoService) GetDetail(ctx context.Context, id int) (domain.Todo, error) {
	detail, err := s.repo.GetDetail(ctx, id)
	if err != nil {
		return domain.Todo{}, err
	}
	if len(detail) == 0 {
		return domain.Todo{}, errors.New(fmt.Sprintf("no row error: %d", id))
	}
	return detail[0], nil
}

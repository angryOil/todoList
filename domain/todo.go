package domain

import (
	"errors"
	"time"
)

type Todo struct {
	Id            int
	UserId        int
	Title         string
	Content       string
	OrderNum      int
	IsDeleted     bool
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}

func CreatedTodo(userId int, title, content string, orderNum int) (Todo, error) {
	if err := validateCreateTodo(title, content, orderNum); err != nil {
		return Todo{}, err
	}

	return Todo{
		UserId:    userId,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		OrderNum:  orderNum,
	}, nil
}

func validateCreateTodo(title, content string, orderNum int) error {
	if title == "" {
		return errors.New("title is empty")
	}
	if content == "" {
		return errors.New("content is empty")
	}
	if orderNum == 0 {
		return errors.New("orderNum is empty")
	}
	return nil
}

func ValidTodoField(todo Todo) error {
	if todo.Id == 0 {
		return errors.New("todoId is zero")
	}
	if todo.UserId == 0 {
		return errors.New("userId is zero")
	}
	if todo.Title == "" {
		return errors.New("title is empty")
	}
	if todo.Content == "" {
		return errors.New("content is empty")
	}
	if todo.OrderNum == 0 {
		return errors.New("orderNum is empty")
	}
	return nil
}

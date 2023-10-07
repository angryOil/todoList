package domain

import (
	"errors"
	"time"
)

type Todo struct {
	Id            int
	Title         string
	Content       string
	OrderNum      int
	IsDeleted     bool
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}

func CreatedTodo(title, content string, orderNum int) (Todo, error) {
	if err := validateCreateTodo(title, content, orderNum); err != nil {
		return Todo{}, err
	}

	return Todo{
		Title:    title,
		Content:  content,
		OrderNum: orderNum,
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

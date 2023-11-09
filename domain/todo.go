package domain

import (
	"errors"
	"time"
)

var _ Todo = (*todo)(nil)

type Todo interface {
	ValidCreate() error
}

const (
	IvnalidUserId = ""
)

type todo struct {
	id            int
	userId        int
	title         string
	content       string
	orderNum      int
	isDeleted     bool
	createdAt     time.Time
	lastUpdatedAt time.Time
}

func (t *todo) ValidCreate() error {
	if t.userId < 1 {
		return errors.New("")
	}
	if t.title == "" {
		return errors.New("")
	}
	return nil
}

//func CreatedTodo(userId int, title, content string, orderNum int) (Todo, error) {
//	if err := validateCreateTodo(title, content, orderNum); err != nil {
//		return Todo{}, err
//	}
//
//	return Todo{
//		UserId:    userId,
//		Title:     title,
//		Content:   content,
//		CreatedAt: time.Now(),
//		OrderNum:  orderNum,
//	}, nil
//}
//
//func validateCreateTodo(title, content string, orderNum int) error {
//	if title == "" {
//		return errors.New("title is empty")
//	}
//	if content == "" {
//		return errors.New("content is empty")
//	}
//	if orderNum == 0 {
//		return errors.New("orderNum is empty")
//	}
//	return nil
//}
//
//func ValidTodoField(todo Todo) error {
//	if todo.Id == 0 {
//		return errors.New("todoId is zero")
//	}
//	if todo.UserId == 0 {
//		return errors.New("userId is zero")
//	}
//	if todo.Title == "" {
//		return errors.New("title is empty")
//	}
//	if todo.Content == "" {
//		return errors.New("content is empty")
//	}
//	if todo.OrderNum == 0 {
//		return errors.New("orderNum is empty")
//	}
//	return nil
//}

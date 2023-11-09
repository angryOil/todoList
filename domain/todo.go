package domain

import (
	"errors"
	"time"
	"todoList/domain/vo"
)

var _ Todo = (*todo)(nil)

type Todo interface {
	ValidCreate() error
	ValidUpdate() error

	Update(title, content string, orderNum int, isDeleted bool) Todo
	ToSave() vo.Save
	ToDetail() vo.Detail
	ToInfo() vo.Info
}

const (
	InvalidUserID   = "invalid user id"
	InvalidTitle    = "invalid title"
	InvalidId       = "invalid id"
	InvalidOrderNum = "invalid order num"
	InvalidContent  = "invalid content"
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

func (t *todo) ToInfo() vo.Info {
	return vo.Info{
		Id:        t.id,
		Title:     t.title,
		UserId:    t.userId,
		CreatedAt: t.createdAt,
		OrderNum:  t.orderNum,
		IsDeleted: t.isDeleted,
	}
}

func (t *todo) ToDetail() vo.Detail {
	return vo.Detail{
		Id:        t.id,
		Title:     t.title,
		UserId:    t.userId,
		Content:   t.content,
		CreatedAt: t.createdAt,
		OrderNum:  t.orderNum,
		IsDeleted: t.isDeleted,
	}
}

func (t *todo) ToSave() vo.Save {
	return vo.Save{
		Id:            t.id,
		UserId:        t.userId,
		Title:         t.title,
		Content:       t.content,
		OrderNum:      t.orderNum,
		IsDeleted:     t.isDeleted,
		CreatedAt:     t.createdAt,
		LastUpdatedAt: t.lastUpdatedAt,
	}
}

func (t *todo) Update(title, content string, orderNum int, isDeleted bool) Todo {
	t.title = title
	t.content = content
	t.orderNum = orderNum
	t.isDeleted = isDeleted
	t.lastUpdatedAt = time.Now()
	return t
}

func (t *todo) ValidUpdate() error {
	if t.id < 1 {
		return errors.New(InvalidId)
	}
	if t.userId < 1 {
		return errors.New(InvalidUserID)
	}
	if t.orderNum < 0 {
		return errors.New(InvalidOrderNum)
	}
	if t.title == "" {
		return errors.New(InvalidTitle)
	}
	if t.content == "" {
		return errors.New(InvalidContent)
	}
	return nil
}

func (t *todo) ValidCreate() error {
	if t.userId < 1 {
		return errors.New(InvalidUserID)
	}
	if t.orderNum < 0 {
		return errors.New(InvalidOrderNum)
	}
	if t.title == "" {
		return errors.New(InvalidTitle)
	}
	if t.content == "" {
		return errors.New(InvalidContent)
	}
	return nil
}

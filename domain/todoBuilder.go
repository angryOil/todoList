package domain

import "time"

var _ TodoBuilder = (*todoBuilder)(nil)

type TodoBuilder interface {
	Id(id int) TodoBuilder
	UserId(userId int) TodoBuilder
	Title(title string) TodoBuilder
	Content(content string) TodoBuilder
	OrderNum(orderNum int) TodoBuilder
	IsDeleted(isDeleted bool) TodoBuilder
	CreatedAt(createdAt time.Time) TodoBuilder
	LastUpdatedAt(lastUpdatedAt time.Time) TodoBuilder
	Build() Todo
}

type todoBuilder struct {
	id            int
	userId        int
	title         string
	content       string
	orderNum      int
	isDeleted     bool
	createdAt     time.Time
	lastUpdatedAt time.Time
}

func (t *todoBuilder) Id(id int) TodoBuilder {
	t.id = id
	return t
}

func (t *todoBuilder) UserId(userId int) TodoBuilder {
	t.userId = userId
	return t
}

func (t *todoBuilder) Title(title string) TodoBuilder {
	t.title = title
	return t
}

func (t *todoBuilder) Content(content string) TodoBuilder {
	t.content = content
	return t
}

func (t *todoBuilder) OrderNum(orderNum int) TodoBuilder {
	t.orderNum = orderNum
	return t
}

func (t *todoBuilder) IsDeleted(isDeleted bool) TodoBuilder {
	t.isDeleted = isDeleted
	return t
}

func (t *todoBuilder) CreatedAt(createdAt time.Time) TodoBuilder {
	t.createdAt = createdAt
	return t
}

func (t *todoBuilder) LastUpdatedAt(lastUpdatedAt time.Time) TodoBuilder {
	t.lastUpdatedAt = lastUpdatedAt
	return t
}

func (t *todoBuilder) Build() Todo {
	return &todo{
		id:            t.id,
		userId:        t.userId,
		title:         t.title,
		content:       t.content,
		orderNum:      t.orderNum,
		isDeleted:     t.isDeleted,
		createdAt:     t.createdAt,
		lastUpdatedAt: t.lastUpdatedAt,
	}
}

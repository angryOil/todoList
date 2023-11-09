package req

import "time"

type CreateTodo struct {
	UserId    int
	Title     string
	Content   string
	OrderNum  int
	IsDeleted bool
	CreatedAt time.Time
}

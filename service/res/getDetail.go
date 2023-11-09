package res

import "time"

type GetDetail struct {
	Id        int
	Title     string
	UserId    int
	Content   string
	CreatedAt time.Time
	OrderNum  int
	IsDeleted bool
}

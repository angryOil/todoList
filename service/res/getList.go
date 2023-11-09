package res

import "time"

type GetList struct {
	Id        int
	UserId    int
	Title     string
	CreatedAt time.Time
	OrderNum  int
	IsDeleted bool
}

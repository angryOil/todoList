package vo

import "time"

type Detail struct {
	Id        int
	Title     string
	UserId    int
	Content   string
	CreatedAt time.Time
	OrderNum  int
	IsDeleted bool
}

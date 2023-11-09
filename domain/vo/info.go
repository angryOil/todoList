package vo

import "time"

type Info struct {
	Id        int
	Title     string
	UserId    int
	CreatedAt time.Time
	OrderNum  int
	IsDeleted bool
}

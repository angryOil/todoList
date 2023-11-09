package req

import "time"

type Save struct {
	Id            int
	UserId        int
	Title         string
	Content       string
	OrderNum      int
	IsDeleted     bool
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}

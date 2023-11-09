package req

type Save struct {
	Id        int
	UserId    int
	Title     string
	Content   string
	OrderNum  int
	IsDeleted bool
}

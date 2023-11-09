package res

type ListDto struct {
	Id        int    `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	OrderNum  int    `json:"order_num,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	IsDeleted bool   `json:"is_deleted,omitempty"`
}

type DetailDto struct {
	Id        int    `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	UserId    int    `json:"user_id,omitempty"`
	Content   string `json:"content,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	OrderNum  int    `json:"order_num,omitempty"`
	IsDeleted bool   `json:"is_deleted,omitempty"`
}

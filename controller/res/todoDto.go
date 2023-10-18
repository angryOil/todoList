package res

import (
	"time"
	"todoList/domain"
)

type ListDto struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	OrderNum  int       `json:"order_num"`
	CreatedAt time.Time `json:"created_at"`
	IsDeleted bool      `json:"is_deleted"`
}

func ToListDtoList(todos []domain.Todo) []ListDto {
	listDto := make([]ListDto, len(todos))
	for i, todo := range todos {
		listDto[i] = ToListDto(todo)
	}
	return listDto
}

func ToListDto(todo domain.Todo) ListDto {
	return ListDto{
		Id:        todo.Id,
		Title:     todo.Title,
		OrderNum:  todo.OrderNum,
		CreatedAt: todo.CreatedAt,
		IsDeleted: todo.IsDeleted,
	}
}

type DetailDto struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	UserId    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	OrderNum  int       `json:"order_num"`
	IsDeleted bool      `json:"is_deleted"`
}

func ToDetailDto(todo domain.Todo) DetailDto {
	return DetailDto{
		Id:        todo.Id,
		UserId:    todo.UserId,
		Title:     todo.Title,
		Content:   todo.Content,
		CreatedAt: todo.CreatedAt,
		OrderNum:  todo.OrderNum,
		IsDeleted: todo.IsDeleted,
	}
}

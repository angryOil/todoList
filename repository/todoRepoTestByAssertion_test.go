package repository

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"todoList/domain"
	"todoList/repository/infla"
)

var repo = NewRepository(
	infla.NewDB(),
)
var ctx = context.Background()
var mockTime = time.Date(2022, 10, 10, 11, 30, 30, 0, utcLocation)

func TestTodoRepository_Create_Assertion(t *testing.T) {
	reqTodo := domain.Todo{
		Id:            9999,
		Title:         "request title",
		Content:       "content 9999",
		OrderNum:      23,
		IsDeleted:     false,
		CreatedAt:     mockTime,
		LastUpdatedAt: mockTime,
	}
	err := repo.Create(ctx, reqTodo)
	assert.NoError(t, err)

	results, err := repo.GetDetail(ctx, reqTodo.Id)
	assert.NoError(t, err)
	assert.NotZero(t, len(results))
	result := results[0]
	assert.Equal(t, reqTodo.Id, result.Id)
	assert.Equal(t, reqTodo.Title, result.Title)
	assert.Equal(t, reqTodo.Content, result.Content)
	assert.Equal(t, reqTodo.OrderNum, result.OrderNum)
	assert.Equal(t, reqTodo.IsDeleted, result.IsDeleted)
	assert.Equal(t, reqTodo.CreatedAt, result.CreatedAt)
	assert.Equal(t, reqTodo.LastUpdatedAt, result.LastUpdatedAt)

	// 정리
	err = repo.Delete(ctx, reqTodo.Id)
	assert.NoError(t, err)

	results, err = repo.GetDetail(ctx, reqTodo.Id)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(results))

}

func TestTodoRepository_Save_Assertion(t *testing.T) {
	firstTodo := domain.Todo{
		Id:            9999,
		Title:         "request title",
		Content:       "content 9999",
		OrderNum:      23,
		IsDeleted:     false,
		CreatedAt:     mockTime,
		LastUpdatedAt: mockTime,
	}
	err := repo.Create(ctx, firstTodo)
	assert.NoError(t, err)

	// 저장 확인
	getTodos, err := repo.GetDetail(ctx, firstTodo.Id)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(getTodos))
	getTodo := getTodos[0]
	assert.Equal(t, firstTodo.Id, getTodo.Id)
	assert.Equal(t, firstTodo.Title, getTodo.Title)
	assert.Equal(t, firstTodo.Content, getTodo.Content)
	assert.Equal(t, firstTodo.OrderNum, getTodo.OrderNum)
	assert.Equal(t, firstTodo.IsDeleted, getTodo.IsDeleted)

	updateTodo := domain.Todo{
		Id:            firstTodo.Id,
		Title:         "Update title",
		Content:       "Update Content ",
		OrderNum:      12,
		IsDeleted:     true,
		CreatedAt:     time.Time{},
		LastUpdatedAt: time.Time{},
	}

	err = repo.Save(ctx, updateTodo, func(td domain.Todo) error {
		if td.Title == "" {
			return errors.New("title is empty")
		}
		return nil
	})
	assert.NoError(t, err)

	getTodos, err = repo.GetDetail(ctx, firstTodo.Id)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(getTodos))
	getTodo = getTodos[0]

	assert.Equal(t, updateTodo.Title, getTodo.Title)
	assert.Equal(t, updateTodo.Content, getTodo.Content)
	assert.Equal(t, updateTodo.OrderNum, getTodo.OrderNum)
	assert.Equal(t, updateTodo.IsDeleted, getTodo.IsDeleted)

	// 정리
	err = repo.Delete(ctx, updateTodo.Id)
	assert.NoError(t, err)
}

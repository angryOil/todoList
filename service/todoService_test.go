package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"todoList/domain"
	"todoList/repository"
	"todoList/repository/infla"
)

var ctx = context.Background()
var service = NewService(
	repository.NewRepository(
		infla.NewDB(),
	),
)

func TestTodoService_CreateTodo(t *testing.T) {
	successSaveDomain := domain.Todo{
		Id:            9999,
		Title:         "만들 타이틀",
		Content:       "만들 컨텐츠",
		OrderNum:      44,
		IsDeleted:     false,
		CreatedAt:     time.Time{},
		LastUpdatedAt: time.Time{},
	}
	err := service.CreateTodo(ctx, successSaveDomain)
	assert.NoError(t, err)
}

func TestTodoService_UpdateTodoSuccess(t *testing.T) {
	successSaveDomain := domain.Todo{
		Id:            9999,
		Title:         "만들 타이틀",
		Content:       "만들 컨텐츠",
		OrderNum:      44,
		IsDeleted:     false,
		CreatedAt:     time.Time{},
		LastUpdatedAt: time.Time{},
	}
	err := service.UpdateTodo(ctx, successSaveDomain)
	assert.NoError(t, err)

	err = service.DeleteTodo(ctx, successSaveDomain.Id)
	assert.NoError(t, err)

	findTodo, err := service.GetDetail(ctx, successSaveDomain.Id)
	assert.Contains(t, err.Error(), "no rows ")
	assert.Equal(t, domain.Todo{}, findTodo)
}

// validate 확인
func TestTodoService_UpdateTodoFail(t *testing.T) {
	noTitleDomain := domain.Todo{
		Id:      9999,
		Content: "만들 컨텐츠",
	}
	err := service.UpdateTodo(ctx, noTitleDomain)
	assert.Contains(t, err.Error(), "title is empty")

	noContentDomain := domain.Todo{
		Id:    9999,
		Title: "제목만있습니다",
	}
	err = service.UpdateTodo(ctx, noContentDomain)
	assert.Contains(t, err.Error(), "content is empty")

	findTodo, err := service.GetDetail(ctx, noTitleDomain.Id)
	assert.Contains(t, err.Error(), "no rows")
	assert.Equal(t, domain.Todo{}, findTodo)
}

// todo 모킹이 되어있지 않기 때문에 데이터가 추가 됨에 따라서
// 자세한 테스트가 불가능 추후 트렌젝션 or dataMocking 사용으로 해결 예정
func TestTodoService_GetTodos(t *testing.T) {

}

func TestTodoService_GetDetail(t *testing.T) {
	successSaveDomain := domain.Todo{
		Id:            7777,
		Title:         "테스트용 title",
		Content:       "content !!",
		OrderNum:      2,
		IsDeleted:     true,
		CreatedAt:     time.Time{},
		LastUpdatedAt: time.Time{},
	}
	err := service.UpdateTodo(ctx, successSaveDomain)
	assert.NoError(t, err)

	result, err := service.GetDetail(ctx, successSaveDomain.Id)
	assert.NoError(t, err)

	assert.Equal(t, successSaveDomain.Id, result.Id)
	assert.Equal(t, successSaveDomain.Title, result.Title)
	assert.Equal(t, successSaveDomain.Content, result.Content)
	assert.Equal(t, successSaveDomain.OrderNum, result.OrderNum)
	assert.Equal(t, successSaveDomain.IsDeleted, result.IsDeleted)
	assert.Equal(t, successSaveDomain.CreatedAt, result.CreatedAt)
	assert.Equal(t, successSaveDomain.LastUpdatedAt, result.LastUpdatedAt)
}

package controller

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"todoList/controller/req"
	"todoList/controller/res"
	"todoList/repository"
	"todoList/repository/infla"
	"todoList/service"
)

var ctx = context.Background()

var c = NewController(
	service.NewService(
		repository.NewRepository(
			infla.NewDB(),
		),
	),
)

func TestTodoController_CreateTodo_Success(t *testing.T) {
	reqDto := req.CreateTodoDto{
		Title:    "hello",
		Content:  "world",
		OrderNum: 2,
	}

	err := c.CreateTodo(ctx, reqDto)
	assert.NoError(t, err)

}

func TestTodoController_CreateTodo_Fail(t *testing.T) {
	reqDto := req.CreateTodoDto{
		Title:    "no content",
		OrderNum: 2,
	}

	err := c.CreateTodo(ctx, reqDto)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "content is empty")
}
func TestTodoController_Update_Success(t *testing.T) {
	reqDto := req.UpdateTodoDto{
		Id:        7777,
		Title:     "ti",
		Content:   "con",
		OrderNum:  2,
		IsDeleted: false,
	}

	err := c.UpdateTodo(ctx, reqDto)
	assert.NoError(t, err)

	result, err := c.GetDetail(ctx, reqDto.Id)
	assert.NoError(t, err)

	assert.Equal(t, result.Id, reqDto.Id)
	assert.Equal(t, result.Title, reqDto.Title)
	assert.Equal(t, result.Content, reqDto.Content)
	assert.Equal(t, result.OrderNum, reqDto.OrderNum)
	assert.Equal(t, result.IsDeleted, reqDto.IsDeleted)

}

func TestTodoController_Update_Fail(t *testing.T) {
	onlyTitleDto := req.UpdateTodoDto{
		Id:    11,
		Title: "no content",
	}

	err := c.UpdateTodo(ctx, onlyTitleDto)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "content is empty")

	noTitleDto := req.UpdateTodoDto{
		Id:        11,
		Content:   "adfs",
		OrderNum:  0,
		IsDeleted: false,
	}
	err = c.UpdateTodo(ctx, noTitleDto)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title is empty")

	noIdDto := req.UpdateTodoDto{
		Title:     "ti",
		Content:   "co",
		OrderNum:  2,
		IsDeleted: true,
	}
	err = c.UpdateTodo(ctx, noIdDto)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is no value")
}

func TestTodoController_GetDetailSuccess(t *testing.T) {
	reqDto := req.UpdateTodoDto{
		Id:        7777,
		Title:     "ti",
		Content:   "con",
		OrderNum:  2,
		IsDeleted: false,
	}

	err := c.UpdateTodo(ctx, reqDto)
	assert.NoError(t, err)

	result, err := c.GetDetail(ctx, reqDto.Id)
	assert.NoError(t, err)

	assert.Equal(t, result.Id, reqDto.Id)
	assert.Equal(t, result.Title, reqDto.Title)
	assert.Equal(t, result.Content, reqDto.Content)
	assert.Equal(t, result.OrderNum, reqDto.OrderNum)
	assert.Equal(t, result.IsDeleted, reqDto.IsDeleted)

}

func TestTodoController_GetDetailFail(t *testing.T) {
	result, err := c.GetDetail(ctx, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is no value")

	assert.Equal(t, res.DetailDto{}, result)

}

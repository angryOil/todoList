package repository

import (
	"context"
	"github.com/uptrace/bun"
	"todoList/domain"
	"todoList/repository/model"
)

type tx func()
type rollback func()
type commit func() error

type transaction interface {
	begin() (error, tx, rollback, commit)
}

type TodoRepository struct {
	db *bun.DB
}

func (r TodoRepository) Create(todo domain.Todo) error {
	return nil
}
func (r TodoRepository) Save(
	todo domain.Todo, saveFunc func(todo2 domain.Todo) error,
) error {
	//err, _, rollback, commit := r.tx.begin()
	//if err != nil {
	//	return err
	//}
	// model := r.db.select(&todo)
	// todoDomain := model.ToDomain
	// err := saveFunc(todoDomain)

	//if err := commit(); err != nil {
	//	fmt.Println("fail to update because of you")
	//	rollback()
	//	return err
	//}
	return nil
}

func (r TodoRepository) GetDetail(ctx context.Context, id int) ([]domain.Todo, error) {
	var result []model.TodoDetail
	err := r.db.NewSelect().Model(&result).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return []domain.Todo{}, err
	}
	return model.ToDomainDetailList(result), nil
}

func (r TodoRepository) GetList(ctx context.Context) ([]domain.Todo, error) {
	var result []model.Todo
	err := r.db.NewSelect().Model(&result).Scan(ctx)
	if err != nil {

		return []domain.Todo{}, err
	}
	return model.ToDomainList(result), nil
	return []domain.Todo{}, nil
}

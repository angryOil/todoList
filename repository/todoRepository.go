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

func (r TodoRepository) Create(ctx context.Context, todo domain.Todo) error {
	tdModel := model.ToDetailModel(todo)
	_, err := r.db.NewInsert().Model(&tdModel).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r TodoRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.NewDelete().Model(&model.Todo{}).Where("id = ?", id).Exec(ctx)
	return err
}

// 사실상 업데이트입니다.
// 있다면 update 있다면 save 입니다 (upsert)

func (r TodoRepository) Save(ctx context.Context, td domain.Todo, saveFunc func(domain.Todo) error) error {
	err := saveFunc(td)
	if err != nil {
		return err
	}
	tdModel := model.ToDetailModel(td)
	_, err = r.db.NewInsert().Model(&tdModel).
		On("CONFLICT (id) DO UPDATE").Exec(ctx)
	return err
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

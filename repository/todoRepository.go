package repository

import (
	"context"
	"github.com/uptrace/bun"
	"todoList/domain"
	"todoList/page"
	"todoList/repository/model"
)

// repository 는 domain 과 model 을 둘다 사용

type tx func()
type rollback func()
type commit func() error

type transaction interface {
	begin() (tx, rollback, commit, error)
}

type TodoRepository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) TodoRepository {
	return TodoRepository{db: db}
}

func (r TodoRepository) GetTransaction(ctx context.Context) (bun.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	return tx, err
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

func (r TodoRepository) Save(ctx context.Context, td domain.Todo, saveValidFunc func(domain.Todo) error) error {
	err := saveValidFunc(td)
	if err != nil {
		return err
	}
	tdModel := model.ToDetailModel(td)
	_, err = r.db.NewInsert().Model(&tdModel).
		On("CONFLICT (id) DO UPDATE").Exec(ctx)
	return err
}
func (r TodoRepository) TxSave(ctx context.Context, tx bun.Tx, td domain.Todo, saveValidFunc func(domain.Todo) error) error {
	err := saveValidFunc(td)
	if err != nil {
		return err
	}
	tdModel := model.ToDetailModel(td)
	_, err = tx.NewInsert().Model(&tdModel).
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

// todo transaction 을 알게 될때 테스트 예정

func (r TodoRepository) GetList(ctx context.Context, page page.ReqPage) ([]domain.Todo, int, error) {
	var result []model.Todo

	// order by desc 는 국룰입니다.
	err := r.db.NewSelect().Model(&result).Limit(page.Size).Offset(page.Page * page.Size).Order("id desc").Scan(ctx)
	if err != nil {

		return []domain.Todo{}, 0, err
	}
	count, err := r.db.NewSelect().Model(&result).Count(ctx)
	return model.ToDomainList(result), count, nil
}

func (r TodoRepository) TxGetList(ctx context.Context, tx bun.Tx, page page.ReqPage) ([]domain.Todo, int, error) {
	var result []model.Todo

	// order by desc 는 국룰입니다.
	err := tx.NewSelect().Model(&result).Limit(page.Size).Offset(page.Page * page.Size).Order("id desc").Scan(ctx)
	if err != nil {

		return []domain.Todo{}, 0, err
	}
	count, err := tx.NewSelect().Model(&result).Count(ctx)
	return model.ToDomainList(result), count, nil
}

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
	db bun.IDB
}

func NewRepository(db bun.IDB) TodoRepository {
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

func (r TodoRepository) Delete(ctx context.Context, userId, id int) error {
	_, err := r.db.NewDelete().Model(&model.Todo{}).Where("id = ? And user_id=?", id, userId).Exec(ctx)
	return err
}

// 사실상 업데이트입니다.
// 있다면 update 있다면 save 입니다 (upsert)

func (r TodoRepository) Save(
	ctx context.Context, userId, id int,
	getValidFunc func([]domain.Todo) (domain.Todo, error),
	mergeTodo func(todo domain.Todo) domain.Todo,
	saveValidFunc func(domain.Todo) error,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	r2 := NewRepository(tx)
	todos, err := r2.GetDetail(ctx, userId, id)
	if err != nil {
		return err
	}

	todo, err := getValidFunc(todos)
	if err != nil {
		return err
	}
	todo = mergeTodo(todo)
	err = saveValidFunc(todo)
	if err != nil {
		return err
	}

	tdModel := model.ToDetailModel(todo)
	_, err = r2.db.NewInsert().Model(&tdModel).
		On("CONFLICT (id) DO UPDATE").Exec(ctx)
	err = tx.Commit()
	return err
}

func (r TodoRepository) GetDetail(ctx context.Context, userId, id int) ([]domain.Todo, error) {
	var result []model.TodoDetail
	err := r.db.NewSelect().Model(&result).Where("id = ? AND user_id = ?", id, userId).Scan(ctx)
	if err != nil {
		return []domain.Todo{}, err
	}
	return model.ToDomainDetailList(result), nil
}

func (r TodoRepository) GetList(ctx context.Context, userId int, page page.ReqPage) ([]domain.Todo, int, error) {
	var result []model.Todo

	// order by desc 는 국룰입니다.
	err := r.db.NewSelect().Model(&result).Where("user_id =?", userId).Limit(page.Size).Offset(page.Page * page.Size).Order("id desc").Scan(ctx)
	if err != nil {

		return []domain.Todo{}, 0, err
	}
	count, err := r.db.NewSelect().Where("user_id=?", userId).Model(&result).Count(ctx)
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

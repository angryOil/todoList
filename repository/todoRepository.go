package repository

import (
	"context"
	"errors"
	"github.com/uptrace/bun"
	"todoList/domain"
	"todoList/domain/vo"
	"todoList/page"
	"todoList/repository/model"
	"todoList/repository/req"
)

// repository 는 domain 과 model 을 둘다 사용

type TodoRepository struct {
	db bun.IDB
}

func NewRepository(db bun.IDB) TodoRepository {
	return TodoRepository{db: db}
}

const (
	InternalServerError = "internal server error"
)

func (r TodoRepository) Create(ctx context.Context, c req.CreateTodo) error {
	tdModel := model.ToCreateModel(c)
	_, err := r.db.NewInsert().Model(&tdModel).Exec(ctx)
	return err
}

func (r TodoRepository) Delete(ctx context.Context, userId, id int) error {
	_, err := r.db.NewDelete().Model(&model.Todo{}).Where("id = ? And user_id=?", id, userId).Exec(ctx)
	return err
}

// 사실상 업데이트입니다.
// 있다면 update 있다면 save 입니다 (upsert)

func (r TodoRepository) Save(
	ctx context.Context, userId, id int,
	getValidFunc func([]domain.Todo) (domain.Todo, error), // 존재하는 데이터 확인
	mergeTodo func(todo domain.Todo) (vo.Save, error), // 저장할 데이터 (기존 데이터와 요청 데이터 병합)
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	r2 := NewRepository(tx)
	todos, err := r2.GetDetail(ctx, userId, id)
	if err != nil {
		return err
	}

	todo, err := getValidFunc(todos)
	if err != nil {
		return err
	}
	v, err := mergeTodo(todo)
	if err != nil {
		return err
	}

	tdModel := model.ToSaveModel(req.Save{
		Id:            v.Id,
		UserId:        v.UserId,
		Title:         v.Title,
		Content:       v.Content,
		OrderNum:      v.OrderNum,
		IsDeleted:     v.IsDeleted,
		CreatedAt:     v.CreatedAt,
		LastUpdatedAt: v.LastUpdatedAt,
	})

	_, err = r2.db.NewInsert().Model(&tdModel).
		On("CONFLICT (id) DO UPDATE").Exec(ctx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return errors.New(InternalServerError)
	}

	return nil
}

func (r TodoRepository) GetDetail(ctx context.Context, userId, id int) ([]domain.Todo, error) {
	var result []model.Todo
	err := r.db.NewSelect().Model(&result).Where("id = ? AND user_id = ?", id, userId).Scan(ctx)
	if err != nil {
		return []domain.Todo{}, err
	}
	return model.ToDomainList(result), nil
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

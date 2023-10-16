package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
	"todoList/domain"
	"todoList/repository/infla"
)

type TodoRepositoryTestSuite struct {
	suite.Suite
	repository    *TodoRepository
	rollback      func() error
	commit        func() error
	koreaLocation *time.Location
}

func TestTodoRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &TodoRepositoryTestSuite{})
}

func (s *TodoRepositoryTestSuite) SetupTest() {
	var db = infla.NewDB()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Panicf("tx setup fail err: %e", err)
	}
	s.rollback = tx.Rollback
	s.commit = tx.Commit
	s.koreaLocation, err = time.LoadLocation("Asia/Seoul")
	if err != nil {
		log.Panicf("location setup fail err: %e", err)
	}

	repository := NewRepository(tx)
	s.repository = &repository
}

func (s *TodoRepositoryTestSuite) AfterTest(suiteName, testName string) {
	log.Printf("roll back / suiteName: %s, testName: %s", suiteName, testName)
	s.commit()
}

func (s *TodoRepositoryTestSuite) TestGetDetail() {
	s.Run("userId 8과 id 16이 주어진다면", func() {
		result, err := s.repository.GetDetail(context.Background(), 8, 16)

		assert.Equal(s.T(), 1, len(result))
		assert.Equal(s.T(), 8, result[0].UserId)
		assert.Equal(s.T(), 16, result[0].Id)
		assert.Nil(s.T(), err)
	})
}

func (s *TodoRepositoryTestSuite) TestSave() {
	s.Run("userId 8이고 id 16인 todo 를 수정할때", func() {
		s.Run("각 request 에 해당하는 값이 모두 주어졌을때", func() {
			userId := 8
			id := 16
			givenDomainTodo := domain.Todo{
				Id:            id,
				UserId:        userId,
				Title:         "mock title",
				Content:       "mock content",
				OrderNum:      100,
				IsDeleted:     false,
				CreatedAt:     time.Date(1, 2, 3, 4, 5, 6, 7, s.koreaLocation),
				LastUpdatedAt: time.Date(1, 2, 3, 4, 5, 6, 7, s.koreaLocation),
			}

			err := s.repository.Save(context.Background(), userId, id,
				func(todos []domain.Todo) (domain.Todo, error) {
					var t = todos[0]
					return domain.Todo{
						Id:            t.Id,
						UserId:        t.UserId,
						Title:         t.Title,
						Content:       t.Content,
						OrderNum:      t.OrderNum,
						IsDeleted:     t.IsDeleted,
						CreatedAt:     t.CreatedAt,
						LastUpdatedAt: t.LastUpdatedAt,
					}, nil
				},
				func(todo domain.Todo) domain.Todo {
					return givenDomainTodo
				},
				func(todo domain.Todo) error {
					return nil
				},
			)

			assert.Nil(s.T(), err)
		})
	})
}

func (s *TodoRepositoryTestSuite) TestSave2() {
	s.Run("userId 8이고 id 16인 todo 를 수정할때", func() {
		s.Run("각 request 에 해당하는 값이 모두 주어졌을때", func() {
			userId := 8
			id := 16
			givenDomainTodo := domain.Todo{
				Id:            id,
				UserId:        userId,
				Title:         "mock title",
				Content:       "mock content",
				OrderNum:      100,
				IsDeleted:     false,
				CreatedAt:     time.Date(1, 2, 3, 4, 5, 6, 7, s.koreaLocation),
				LastUpdatedAt: time.Date(1, 2, 3, 4, 5, 6, 7, s.koreaLocation),
			}

			err := s.repository.Save(context.Background(), userId, id,
				func(todos []domain.Todo) (domain.Todo, error) {
					var t = todos[0]
					return domain.Todo{
						Id:            t.Id,
						UserId:        t.UserId,
						Title:         t.Title,
						Content:       t.Content,
						OrderNum:      t.OrderNum,
						IsDeleted:     t.IsDeleted,
						CreatedAt:     t.CreatedAt,
						LastUpdatedAt: t.LastUpdatedAt,
					}, nil
				},
				func(todo domain.Todo) domain.Todo {
					return givenDomainTodo
				},
				func(todo domain.Todo) error {
					return nil
				},
			)

			assert.Nil(s.T(), err)
		})
	})
}

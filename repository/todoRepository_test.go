package repository

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"todoList/domain"
	"todoList/page"
	"todoList/repository/infla"
)

type TodoRepositoryTestSuite struct {
	suite.Suite
	repository    ITodoRepository
	rollback      func() error
	commit        func() error
	koreaLocation *time.Location
}

func TestTodoRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &TodoRepositoryTestSuite{})
}

// 각 테스트실행전 1번씩 실행됨 전체 테스트 실해전 실행은 SetupSuite() 메소드 사용
func (s *TodoRepositoryTestSuite) SetupTest() {
	log.Printf("setup testSuite ...")
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
	s.repository = repository
}

// userId 가9999 인 테스트 데이터를 10개 가지고 시작할거임

// 테스트시 기준이 될 데이터
var baseTestDomain = domain.Todo{
	UserId:    9999,
	Title:     "base test model title",
	Content:   "base test model content",
	OrderNum:  1,
	IsDeleted: false,
}

func getInitDomainArr() []domain.Todo {
	base := baseTestDomain
	result := make([]domain.Todo, 10)
	for i, _ := range result {
		result[i] = domain.Todo{
			Id:        i + 1, //id 가 0이 되면안되므로
			UserId:    base.UserId,
			Title:     base.Title + strconv.Itoa(i),
			Content:   base.Content + strconv.Itoa(i),
			OrderNum:  base.OrderNum + i,
			IsDeleted: base.IsDeleted,
		}
	}
	return result
}

// test 에선 멱등성을 위해 save(update) 메소드를 사용할거고
// update 에서는 getValidFunc 에서 나온 결과로(repo 에서 조회후 넣을 예정)
func validFunMaker(todoId int) func([]domain.Todo) (domain.Todo, error) {
	return func(todos []domain.Todo) (domain.Todo, error) {
		return domain.Todo{
			UserId:    baseTestDomain.UserId,
			Id:        todoId,
			Title:     baseTestDomain.Title + strconv.Itoa(todoId),
			Content:   baseTestDomain.Content + strconv.Itoa(todoId),
			OrderNum:  baseTestDomain.OrderNum + todoId,
			IsDeleted: baseTestDomain.IsDeleted,
		}, nil
	}
}

// 병합된 데이터 테스트 이므로 그대로
func getTestMergeTodo(todo domain.Todo) domain.Todo {
	return todo
}

// todoId,userId 는 0이어서는 안됨
func getTestSaveValid(todo domain.Todo) error {
	if todo.UserId == 0 {
		return errors.New("userId is zero")
	}
	if todo.Id == 0 {
		return errors.New("todoId is zero")
	}
	return nil
}

var (
	Repo = "repo"
)

var RepoTestFun = func() {
	log.Println("test")
}

// 각테스트전에 실행
func (s *TodoRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	log.Printf("이것은 BeforeTest 입니다. %s %s", suiteName, testName)
	getTestValidFunc := validFunMaker
	for _, todoDomain := range getInitDomainArr() {
		err := s.repository.Save(context.Background(), todoDomain.UserId, todoDomain.Id, getTestValidFunc(todoDomain.Id), getTestMergeTodo, getTestSaveValid)
		if err != nil {
			panic(err)
		}
	}
	results, totalCnt, _ := s.repository.GetList(context.Background(), baseTestDomain.UserId, page.NewReqPage(0, 0))
	log.Println(results, totalCnt)
}

// 물론 BeforeTest 도 존재함 (suiteName,testName string)을 인자로 받아서 실행
func (s *TodoRepositoryTestSuite) AfterTest(suiteName, testName string) {
	log.Printf("roll back / suiteName: %s, testName: %s", suiteName, testName)
	//err := s.commit()
	err := s.rollback()
	if err != nil {
		panic(err)
	}
}

func (s *TodoRepositoryTestSuite) TestGetDetail() {
	targetDomain := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)]
	log.Printf("TestGetDetail target id is %d", targetDomain.Id)
	s.Run("테스트 객체들중 하나를 가져온다", func() {
		result, err := s.repository.GetDetail(context.Background(), targetDomain.UserId, targetDomain.Id)

		// assertion 으로 처리도 가능하지만 suite 를 사용해서도 검증할수있음 선호도에 따라 사용
		s.Equal(1, len(result)) // 이코드는 아래 코드와 같은 역활을함
		assert.Equal(s.T(), 1, len(result))
		assert.Equal(s.T(), targetDomain.UserId, result[0].UserId)
		assert.Equal(s.T(), targetDomain.Id, result[0].Id)
		assert.Nil(s.T(), err)
	})
}

func (s *TodoRepositoryTestSuite) TestSave() {
	s.Run("save(update)시 값변경되는지 확인", func() {
		userId := baseTestDomain.UserId
		id := getInitDomainArr()[0].Id
		getTestValidFunc := validFunMaker

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
			getTestValidFunc(userId),
			func(todo domain.Todo) domain.Todo {
				return givenDomainTodo
			},
			func(todo domain.Todo) error {
				return nil
			},
		)
		assert.Nil(s.T(), err)
	})
}

func (s *TodoRepositoryTestSuite) TestSaveFail() {
	s.Run("save(update)인자들에서 error 가나올경우 error를 반환하는지 확인", func() {
		userId := baseTestDomain.UserId
		id := getInitDomainArr()[0].Id
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
				return domain.Todo{}, errors.New("validFunc fail")
			},
			func(todo domain.Todo) domain.Todo {
				return givenDomainTodo
			},
			func(todo domain.Todo) error {
				return nil
			},
		)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), "validFunc fail", err.Error())

		getValidFunc := validFunMaker
		err = s.repository.Save(context.Background(), userId, id,
			getValidFunc(userId),
			func(todo domain.Todo) domain.Todo {
				return domain.Todo{}
			},
			func(todo domain.Todo) error {
				if todo.UserId == 0 {
					return errors.New("userId is zero")
				}
				return nil
			},
		)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), "userId is zero", err.Error())
	})
}

func (s *TodoRepositoryTestSuite) TestCreate() {
	s.Run("특정 id값의 todoTotal cnt를 가져온후 추가하면 갯수가 늘어있는지 확인한다", func() {
		basePage := page.NewReqPage(0, 0)
		// 추가전 특정id의 todoList의 총갯수
		_, beforeTotalCnt, err := s.repository.GetList(context.Background(), baseTestDomain.UserId, basePage)
		assert.Nil(s.T(), err)
		// todo추가
		err = s.repository.Create(context.Background(), baseTestDomain)
		assert.Nil(s.T(), err)

		// 저장후
		results, afterTotalCnt, err := s.repository.GetList(context.Background(), baseTestDomain.UserId, basePage)
		result := results[0]
		assert.Nil(s.T(), err)
		assert.Equal(s.T(), beforeTotalCnt+1, afterTotalCnt)
		assert.Equal(s.T(), baseTestDomain.Title, result.Title)
		assert.Equal(s.T(), baseTestDomain.OrderNum, result.OrderNum)
		assert.Equal(s.T(), baseTestDomain.IsDeleted, result.IsDeleted)
	})
}

func (s *TodoRepositoryTestSuite) TestGetList() {
	s.Run("저장된 값이 올바른지 확인", func() {
		basePage := page.NewReqPage(0, 0)
		results, _, err := s.repository.GetList(context.Background(), baseTestDomain.UserId, basePage)
		assert.Nil(s.T(), err)
		for _, v := range results {
			assert.NotEqual(s.T(), 0, v.Id)
			assert.NotEqual(s.T(), 0, v.UserId)
			assert.NotEqual(s.T(), "", v.Title)
			assert.NotEqual(s.T(), 0, v.OrderNum)
		}
	})
}

func (s *TodoRepositoryTestSuite) TestDelete() {
	s.Run("삭제후 해당 content가없고 전체 갯수가 감소했는지 확인", func() {
		targetTodo := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)]

		// 존재하는지 확인
		results, err := s.repository.GetDetail(context.Background(), targetTodo.UserId, targetTodo.Id)
		assert.NoError(s.T(), err)
		assert.NotZero(s.T(), len(results))
		assert.Equal(s.T(), targetTodo.Id, results[0].Id)
		assert.NotEqual(s.T(), 0, results[0].Id)

		// 삭제
		err = s.repository.Delete(context.Background(), targetTodo.UserId, targetTodo.Id)
		assert.NoError(s.T(), err)

		// 삭제후 삭제 됐는지 확인
		results, err = s.repository.GetDetail(context.Background(), targetTodo.UserId, targetTodo.Id)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 0, len(results))
	})
}

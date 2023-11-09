package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"todoList/domain"
	"todoList/domain/vo"
	"todoList/page"
	"todoList/repository"
	"todoList/repository/infla"
	"todoList/service/req"
	"todoList/service/res"
)

type TodoServiceTestSuite struct {
	suite.Suite
	service  ITodoService
	rollback func() error
}

func TestTodoServiceTestSuite(t *testing.T) {
	suite.Run(t, &TodoServiceTestSuite{})
}

// 테스트시 기준이 될 데이터
var baseVo = vo.Detail{
	UserId:    9999,
	Title:     "base test model title",
	Content:   "base test model content",
	OrderNum:  1,
	IsDeleted: false,
}

var baseTestPageReq = page.NewReqPage(0, 0)

func getInitDomainArr() []domain.Todo {
	base := baseVo
	result := make([]domain.Todo, 10)
	for i, _ := range result {
		id := i + 1
		result[i] = domain.NewTodoBuilder().
			Id(id).
			UserId(base.UserId).
			Title(base.Title + strconv.Itoa(id)).
			Content(base.Content + strconv.Itoa(id)).
			OrderNum(base.OrderNum + id).
			IsDeleted(base.IsDeleted).
			Build()
	}
	return result
}

// suite 안에 repo 를 넣지 않기위해 따로 만듬 SetupTest, BeforeTest , AfterTest 메서드에서만 사용될거임
var testRepo repository.ITodoRepository

func (s *TodoServiceTestSuite) SetupTest() {
	log.Println("setup service testSuite ...")

	var db = infla.NewDB()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	s.rollback = tx.Rollback

	testRepo = repository.NewRepository(tx)

	s.service = NewService(testRepo)
}

// 병합된 데이터 테스트 이므로 그대로
func getTestMergeTodo(todo domain.Todo) (vo.Save, error) {
	return todo.ToSave(), nil
}

// test 에선 멱등성을 위해 save(update) 메소드를 사용할거고
// update 에서는 getValidFunc 에서 나온 결과로(repo 에서 조회후 넣을 예정)
func validFunMaker(todoId int) func([]domain.Todo) (domain.Todo, error) {
	return func(todos []domain.Todo) (domain.Todo, error) {
		return domain.NewTodoBuilder().
			UserId(baseVo.UserId).
			Id(todoId).
			Title(baseVo.Title + strconv.Itoa(todoId)).
			Content(baseVo.Content + strconv.Itoa(todoId)).
			OrderNum(baseVo.OrderNum + todoId).
			IsDeleted(baseVo.IsDeleted).
			Build(), nil
	}
}

// todoId,userId 는 0이어서는 안됨
func getTestSaveValid(todo domain.Todo) error {
	return todo.ValidUpdate()
}

func (s *TodoServiceTestSuite) BeforeTest(suiteName, testName string) {
	log.Printf("이것은 BeforeTest 입니다. 기본 값들을 생성하죠. %s %s", suiteName, testName)
	getTestValidFunc := validFunMaker
	for _, todoDomain := range getInitDomainArr() {
		v := todoDomain.ToDetail()
		err := testRepo.Save(context.Background(), v.UserId, v.Id, getTestValidFunc(v.Id), getTestMergeTodo)
		if err != nil {
			panic(err)
		}
	}
}
func (s *TodoServiceTestSuite) AfterTest(suiteName, testName string) {
	log.Printf("이것은 AfterTest 입니다. 롤백처리를 하죠. %s %s", suiteName, testName)
	err := s.rollback()
	if err != nil {
		panic(err)
	}
}

func (s *TodoServiceTestSuite) TestCreateSuccess() {
	s.Run("todoList 갯수를 확인후 생성 , 생성후 갯수가 증가했는지 확인", func() {
		// 기존 갯수
		_, beforeCnt, err := s.service.GetTodos(context.Background(), baseVo.UserId, baseTestPageReq)
		assert.NoError(s.T(), err)

		// 생성
		err = s.service.CreateTodo(context.Background(), req.CreateTodo{
			UserId:   baseVo.UserId,
			Title:    baseVo.Title,
			Content:  baseVo.Content,
			OrderNum: baseVo.OrderNum,
		})
		assert.NoError(s.T(), err)

		// 확인
		_, afterCnt, err := s.service.GetTodos(context.Background(), baseVo.UserId, baseTestPageReq)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), afterCnt, beforeCnt+1)
	})
	s.Run("userId에 0을 넣었을경우", func() {
		_, totalCnt, err := s.service.GetTodos(context.Background(), 0, baseTestPageReq)
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "invalid")
		assert.Equal(s.T(), 0, totalCnt)
	})
}
func (s *TodoServiceTestSuite) TestCreateFail() {
	s.Run("제목이 없을경우 제목이없다는 에러를 리턴 ", func() {
		reqTodo := baseVo
		reqTodo.Title = ""
		err := s.service.CreateTodo(context.Background(), req.CreateTodo{
			UserId:   reqTodo.UserId,
			Title:    reqTodo.Title,
			Content:  reqTodo.Content,
			OrderNum: reqTodo.OrderNum,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "invalid")
	})

	s.Run("컨텐츠가 없을경우 컨텐츠가 없다는 에러를 리턴", func() {
		reqTodo := baseVo
		reqTodo.Content = ""
		err := s.service.CreateTodo(context.Background(), req.CreateTodo{
			UserId:   reqTodo.UserId,
			Title:    reqTodo.Title,
			Content:  reqTodo.Content,
			OrderNum: reqTodo.OrderNum,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "invalid")
	})
	s.Run("orderNum 이 0 일경우 에러를 리턴 ", func() {
		reqTodo := baseVo
		reqTodo.OrderNum = 0
		err := s.service.CreateTodo(context.Background(), req.CreateTodo{
			UserId:   reqTodo.UserId,
			Title:    reqTodo.Title,
			Content:  reqTodo.Content,
			OrderNum: reqTodo.OrderNum,
		})
		assert.NoError(s.T(), err)
	})
}

func (s *TodoServiceTestSuite) TestDelete() {
	s.Run("삭제후 정말 갯수가 감소하는지 확인", func() {
		// 기존값 확인
		results, beforeCnt, err := s.service.GetTodos(context.Background(), baseVo.UserId, baseTestPageReq)
		assert.NoError(s.T(), err)
		assert.NotZero(s.T(), len(results))
		targetTodo := results[rand.Intn(len(results)-1)]

		// 삭제
		err = s.service.DeleteTodo(context.Background(), targetTodo.UserId, targetTodo.Id)
		assert.NoError(s.T(), err)

		// 삭제후 값 확인
		_, after, err := s.service.GetTodos(context.Background(), baseVo.UserId, baseTestPageReq)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), beforeCnt-1, after)

		// 없는값을 다시 삭제할경우
		err = s.service.DeleteTodo(context.Background(), targetTodo.UserId, targetTodo.Id)
		assert.NoError(s.T(), err)
		_, after2, err := s.service.GetTodos(context.Background(), baseVo.UserId, baseTestPageReq)
		assert.Equal(s.T(), after, after2)
	})
}

func (s *TodoServiceTestSuite) TestUpdate() {
	s.Run("update 메소드 호출후 변경되었는지 확인", func() {
		// 저장되어있는 값중하나를 랜덤으로 선택(BeforeTest 에 저장함)
		target := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)].ToDetail()

		result, err := s.service.GetDetail(context.Background(), target.UserId, target.Id)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), target.Id, result.Id)
		assert.Equal(s.T(), target.UserId, result.UserId)
		assert.Equal(s.T(), target.Title, result.Title)
		assert.Equal(s.T(), target.Content, result.Content)

		modifyDomain := domain.NewTodoBuilder().
			Id(target.Id).
			UserId(target.UserId).
			Title("mod title").
			Content("mod content").
			OrderNum(22).
			IsDeleted(true).
			Build().ToDetail()

		// 업데이트
		err = s.service.UpdateTodo(context.Background(), req.Save{
			Id:        modifyDomain.Id,
			UserId:    modifyDomain.UserId,
			Title:     modifyDomain.Title,
			Content:   modifyDomain.Content,
			OrderNum:  modifyDomain.OrderNum,
			IsDeleted: modifyDomain.IsDeleted,
		})
		assert.NoError(s.T(), err)

		// 수정한 값과 같은지 확인
		result, err = s.service.GetDetail(context.Background(), modifyDomain.UserId, modifyDomain.Id)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), modifyDomain.Id, result.Id)
		assert.Equal(s.T(), modifyDomain.UserId, result.UserId)
		assert.Equal(s.T(), modifyDomain.OrderNum, result.OrderNum)
		assert.Equal(s.T(), modifyDomain.Title, result.Title)
		assert.Equal(s.T(), modifyDomain.Content, result.Content)

		assert.Equal(s.T(), "mod title", result.Title)
		assert.Equal(s.T(), "mod content", result.Content)
	})
	s.Run("update 메소드에 잘못된 값을 주었을 경우", func() {
		target := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)].ToDetail()
		// 빈값이었을경우
		err := s.service.UpdateTodo(context.Background(), req.Save{})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "no rows")

		// id가 없을경우
		err = s.service.UpdateTodo(context.Background(), req.Save{
			UserId:    target.UserId,
			Title:     "mod title",
			Content:   "mod content",
			OrderNum:  target.OrderNum,
			IsDeleted: target.IsDeleted,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "no rows")

		// userId 가 없을경
		err = s.service.UpdateTodo(context.Background(), req.Save{
			Id:        target.Id,
			Title:     "mod title",
			Content:   "mod content",
			OrderNum:  target.OrderNum,
			IsDeleted: target.IsDeleted,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "no rows")
	})
	s.Run("없는 todo를 수정하려 할경우 error를 반환한다", func() {
		target := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)].ToDetail()
		// 값부터 삭제
		err := s.service.DeleteTodo(context.Background(), target.UserId, target.Id)
		assert.NoError(s.T(), err)
		// 변경 시도
		err = s.service.UpdateTodo(context.Background(), req.Save{
			Id:        target.Id,
			UserId:    target.UserId,
			Title:     target.Title,
			Content:   target.Content,
			OrderNum:  target.OrderNum,
			IsDeleted: target.IsDeleted,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "no row")
	})
}

func (s *TodoServiceTestSuite) TestGetDetail() {
	s.Run("id가 0일경우 error를 반환한다", func() {
		ctx := context.Background()
		_, err := s.service.GetDetail(ctx, 10, 0)
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "invalid")
	})
	s.Run("결과 row 가없는 경우 error를 반환한다", func() {
		// 값을 삭제 (row 가 없어짐)
		ctx := context.Background()
		ctx = context.WithValue(ctx, "userId", baseVo.UserId)
		target := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)].ToDetail()
		err := s.service.DeleteTodo(ctx, target.UserId, target.Id)
		assert.NoError(s.T(), err)

		// 조회
		todo, err := s.service.GetDetail(ctx, target.UserId, target.Id)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), res.GetDetail{}, todo)
	})
}

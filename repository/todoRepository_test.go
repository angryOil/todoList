package repository

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"reflect"
	"testing"
	"time"
	"todoList/domain"
)

var utcLocation, _ = time.LoadLocation("")

func TestTodoRepository_GetList(t *testing.T) {
	dsn := "postgres://postgres:@localhost:5432/postgres?sslmode=disable"
	db := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn))), pgdialect.New())
	worngDb := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN("postgres://postgres:@localhost:486/postgres?sslmode=disable"))), pgdialect.New())
	utcLocation, _ := time.LoadLocation("")

	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Todo
		wantErr bool
	}{
		{
			"select id is 1",
			fields{db: db},
			args{
				ctx: context.Background(),
			},
			[]domain.Todo{
				{
					Id:            1,
					Title:         "인생은 쓰다 하..",
					OrderNum:      1,
					IsDeleted:     false,
					CreatedAt:     time.Date(2022, 10, 10, 11, 30, 30, 0, utcLocation),
					LastUpdatedAt: time.Date(2022, 10, 10, 11, 30, 30, 0, utcLocation),
				},
			},
			false,
		},
		{
			"if db error ",
			fields{db: worngDb},
			args{
				ctx: context.Background(),
			},
			[]domain.Todo{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := TodoRepository{
				db: tt.fields.db,
			}
			got, err := r.GetList(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodoRepository_GetDetail(t *testing.T) {
	dsn := "postgres://postgres:@localhost:5432/postgres?sslmode=disable"
	db := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn))), pgdialect.New())
	wrongDb := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN("postgres://postgres:@localhost:486/postgres?sslmode=disable"))), pgdialect.New())

	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Todo
		wantErr bool
	}{
		{
			"select id is 1",
			fields{db: db},
			args{
				ctx: context.Background(),
				id:  1,
			},
			[]domain.Todo{
				{
					Id:            1,
					Title:         "인생은 쓰다 하..",
					Content:       "오늘 집에 오다가 지하철에 지갑을 두고 내렸다",
					OrderNum:      1,
					IsDeleted:     false,
					CreatedAt:     time.Date(2022, 10, 10, 11, 30, 30, 0, utcLocation),
					LastUpdatedAt: time.Date(2022, 10, 10, 11, 30, 30, 0, utcLocation),
				},
			},
			false,
		},
		{
			"select id is 2 will be error",
			fields{db: wrongDb},
			args{
				ctx: context.Background(),
				id:  2,
			},
			[]domain.Todo{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := TodoRepository{
				db: tt.fields.db,
			}
			got, err := r.GetDetail(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDetail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodoRepository_Create(t *testing.T) {
	dsn := "postgres://postgres:@localhost:5432/postgres?sslmode=disable"
	db := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn))), pgdialect.New())

	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx  context.Context
		todo domain.Todo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "right create test",
			fields: fields{db: db},
			args: args{
				ctx: context.Background(),
				todo: domain.Todo{
					Id:            0,
					Title:         "새로운 타이틀",
					Content:       "새로운 글",
					OrderNum:      2,
					IsDeleted:     false,
					CreatedAt:     time.Date(2022, 10, 10, 11, 30, 30, 0, utcLocation),
					LastUpdatedAt: time.Date(2022, 10, 10, 11, 30, 30, 0, utcLocation),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := TodoRepository{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.todo); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

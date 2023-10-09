package infla

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var dsn = "postgres://postgres:@localhost:5432/postgres?sslmode=disable"
var wrongDsn = "postgres://postgres:@localhost:54321/postgres?sslmode=disable"
var wrongDB = bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(wrongDsn))), pgdialect.New())
var db = bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn))), pgdialect.New())

func NewDB() *bun.DB {
	return db
}

func WrongDB() *bun.DB {
	return wrongDB
}

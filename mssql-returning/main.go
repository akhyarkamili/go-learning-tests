package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
)

func main() {

}

type mssql struct {
	db *sqlx.DB
}

func newMssql(user, pass, host, dbname string) (*mssql, error) {
	db, err := sqlx.Connect("sqlserver", fmt.Sprintf(
		"sqlserver://%s:%s@%s?database=%s&encrypt=disable",
		user, pass, host, dbname,
	))
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	return &mssql{
		db: db,
	}, nil
}

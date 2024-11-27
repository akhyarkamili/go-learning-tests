package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newMssql(user, pass, host, dbname string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlserver", fmt.Sprintf(
		"sqlserver://%s:%s@%s?database=%s&encrypt=disable",
		user, pass, host, dbname,
	))
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	return db, nil
}

func Test_newMssql(t *testing.T) {
	db, err := newMssql("sa", "yourStrong(!)Password", "localhost", "test")
	if err != nil {
		t.Fatalf("newMssql() error = %v", err)
	}

	db.MustExec("create schema s1")
	db.MustExec(`
		create table s1.t1 (
			id int primary key,
			last_update datetime,
		);`)
	t.Cleanup(func() {
		db.MustExec("drop table s1.t1")
		db.MustExec("drop schema s1")
	})

	// insert works
	db.MustExec("INSERT INTO s1.t1 (id, last_update) VALUES (1, '2000-01-01T15:04:05Z')")

	assert.Equal(t,
		[]rowData{
			{"id": int64(1), "last_update": time.Date(2000, 1, 1, 15, 4, 5, 0, time.UTC)},
		},
		getAllData(t, db, "s1.t1", "id"),
	)

	// update returning
	rows, err := db.Queryx(`
		UPDATE s1.t1 SET last_update = '2024-01-01T15:04:05Z' OUTPUT INSERTED.id as inserted_id, DELETED.id AS deleted_id, DELETED.last_update AS deleted_last_update, INSERTED.last_update AS inserted_last_update

	`)

	if !assert.NoError(t, err) {
		assert.FailNow(t, "error reading db data")
	}

	var ret []rowData
	for rows.Next() {
		row := rowData{}
		if !assert.NoError(t, rows.MapScan(row)) {
			assert.FailNow(t, "error scanning row")
		}

		for k, v := range row {
			if v, ok := v.(time.Time); ok {
				row[k] = v.UTC()
			}
		}

		ret = append(ret, row)
	}

	ret[0]["inserted_last_update"] = ret[0]["inserted_last_update"].(time.Time).UTC()
	ret[0]["deleted_last_update"] = ret[0]["deleted_last_update"].(time.Time).UTC()
	// id is correct, unchanged
	assert.Equal(t,
		rowData{
			"inserted_id":          int64(1),
			"deleted_id":           int64(1),
			"inserted_last_update": time.Date(2024, 1, 1, 15, 4, 5, 0, time.UTC),
			"deleted_last_update":  time.Date(2000, 1, 1, 15, 4, 5, 0, time.UTC),
		},
		ret[0],
	)
}

type rowData map[string]any

func getAllData(t *testing.T, db *sqlx.DB, table string, order string) []rowData {
	t.Helper()

	rows, err := db.Queryx(fmt.Sprintf(
		"select * from %s order by %s",
		table, order,
	))
	if !assert.NoError(t, err) {
		assert.FailNow(t, "error reading db data")
	}

	var ret []rowData
	for rows.Next() {
		row := rowData{}
		if !assert.NoError(t, rows.MapScan(row)) {
			assert.FailNow(t, "error scanning row")
		}

		for k, v := range row {
			if v, ok := v.(time.Time); ok {
				row[k] = v.UTC()
			}
		}

		ret = append(ret, row)
	}

	return ret
}

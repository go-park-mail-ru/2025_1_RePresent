package adv

import (
	"database/sql"
)

type CassandraSession interface {
	Query(stmt string, values ...interface{}) Query
	Close()
}

type Query interface {
	Scan(dest ...interface{}) error
	Exec() error
	Iter() Iter
}

type Iter interface {
	Scan(dest ...interface{}) bool
	Close() error
}

type ClickHouseDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

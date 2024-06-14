package sqlx

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func PrepareQueries(db *sqlx.DB, queries []string) ([]*sqlx.Stmt, error) {
	statements := make([]*sqlx.Stmt, len(queries))
	for idx, query := range queries {
		stmt, err := db.Preparex(query)
		if err != nil {
			return nil, err
		}
		statements[idx] = stmt
	}
	return statements, nil
}

func PrepareNamedQueries(db *sqlx.DB, queries []string) ([]*sqlx.NamedStmt, error) {
	statements := make([]*sqlx.NamedStmt, len(queries))
	for idx, query := range queries {
		stmt, err := db.PrepareNamed(query)
		if err != nil {
			return nil, err
		}
		statements[idx] = stmt
	}
	return statements, nil
}

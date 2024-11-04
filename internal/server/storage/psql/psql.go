package psql

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func PingDB(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return err
	}
	return nil
}

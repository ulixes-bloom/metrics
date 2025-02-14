package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	appErrors "github.com/ulixes-bloom/ya-metrics/internal/pkg/errors"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type pgstorage struct {
	db *sql.DB
}

func NewStorage(ctx context.Context, db *sql.DB) (*pgstorage, error) {
	newStorage := pgstorage{db: db}

	if err := newStorage.createTables(ctx); err != nil {
		return nil, fmt.Errorf("pg.NewStorage: %w", err)
	}

	if err := newStorage.PingDB(); err != nil {
		return nil, fmt.Errorf("pg.NewStorage: %w", err)
	}

	return &newStorage, nil
}

func (ps *pgstorage) createTables(ctx context.Context) error {
	_, err := ps.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS metrics
		(
			id    varchar(255) PRIMARY KEY, 
			type  varchar(30) NOT NULL, 
			delta bigint, 
			value double precision
		);`)
	if err != nil {
		return fmt.Errorf("pg.createTables.metrics: %w", err)
	}
	return nil
}

func (ps *pgstorage) Shutdown(ctx context.Context) error {
	return ps.db.Close()
}

func (ps *pgstorage) PingDB() error {
	if err := ps.db.Ping(); err != nil {
		return fmt.Errorf("pg.pingDB: %w", err)
	}
	return nil
}

func (ps *pgstorage) Set(ctx context.Context, metric metrics.Metric) (metrics.Metric, error) {
	if metric.MType != metrics.Counter && metric.MType != metrics.Gauge {
		return metric, appErrors.ErrMetricTypeNotImplemented
	}

	_, err := ps.db.ExecContext(ctx, `
		INSERT INTO metrics (id, type, delta, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) 
		DO UPDATE SET id=$1, type=$2, delta=metrics.delta+$3, value=$4`, metric.ID, metric.MType, metric.Delta, metric.Value)
	if err != nil {
		return metric, fmt.Errorf("pg.set: %w", err)
	}
	return metric, nil
}

func (ps *pgstorage) SetAll(ctx context.Context, meticsSlice []metrics.Metric) error {
	if len(meticsSlice) == 0 {
		return nil
	}

	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("pg.setAll.begin: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO metrics (id, type, delta, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id)
		DO UPDATE SET id=$1, type=$2, delta=metrics.delta+$3, value=$4`)
	if err != nil {
		return fmt.Errorf("pg.setAll.stmtPrepare: %w", err)
	}

	for _, m := range meticsSlice {
		if m.MType != metrics.Counter && m.MType != metrics.Gauge {
			return appErrors.ErrMetricTypeNotImplemented
		}
		_, err = stmt.ExecContext(ctx, m.ID, m.MType, m.Delta, m.Value)
		if err != nil {
			return fmt.Errorf("pg.setAll.stmtExec: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("pg.setAll.commit: %w", err)
	}

	return nil
}

func (ps *pgstorage) Get(ctx context.Context, name string) (val metrics.Metric, ok error) {
	var metric metrics.Metric
	row := ps.db.QueryRowContext(ctx, `
		SELECT id, type, delta, value
		FROM metrics
		WHERE id=$1`, name)
	if err := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value); err != nil {
		return metric, fmt.Errorf("pg.get: %w", err)
	}
	return metric, nil
}

func (ps *pgstorage) GetAll(ctx context.Context) ([]metrics.Metric, error) {
	rows, err := ps.db.QueryContext(ctx, `
		SELECT id, type, delta, value
		FROM metrics`)
	if err != nil {
		return nil, fmt.Errorf("pg.getAll.query: %w", err)
	}
	defer rows.Close()

	allMetrics := []metrics.Metric{}
	for rows.Next() {
		var metric metrics.Metric
		if err := rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value); err != nil {
			return nil, fmt.Errorf("pg.getAll.rowsScan: %w", err)
		}
		allMetrics = append(allMetrics, metric)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("pg.getAll.rowsErr: %w", err)
	}
	return allMetrics, nil
}

func PingDB(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("pg.pingDB.openSql: %w", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return fmt.Errorf("pg.pingDB.ping: %w", err)
	}
	return nil
}

func NeedToRetry() func(err error) bool {
	return func(err error) bool {
		var pgErr *pgconn.PgError
		return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ConnectionException
	}
}

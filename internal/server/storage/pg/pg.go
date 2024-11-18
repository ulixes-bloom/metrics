package pg

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metricerrors"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type pgstorage struct {
	db *sql.DB
}

func NewStorage(dsn string) (*pgstorage, error) {
	newStorage := pgstorage{}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	newStorage.db = db

	err = newStorage.PingDB()
	if err != nil {
		return nil, err
	}

	err = newStorage.Setup()
	if err != nil {
		return nil, err
	}

	return &newStorage, nil
}

func (ps *pgstorage) Setup() error {
	_, err := ps.db.Exec(createTableQuery)
	if err != nil {
		return err
	}
	return nil
}

func (ps *pgstorage) Shutdown() error {
	return ps.db.Close()
}

func (ps *pgstorage) PingDB() error {
	if err := ps.db.Ping(); err != nil {
		return err
	}
	return nil
}

func (ps *pgstorage) Set(metric metrics.Metric) (metrics.Metric, error) {
	if metric.MType != metrics.Counter && metric.MType != metrics.Gauge {
		return metric, metricerrors.ErrMetricTypeNotImplemented
	}

	_, err := ps.db.Exec(setMetricQuery, metric.ID, metric.MType, metric.Delta, metric.Value)
	if err != nil {
		return metric, err
	}
	return metric, nil
}

func (ps *pgstorage) SetAll(meticsSlice []metrics.Metric) error {
	if len(meticsSlice) == 0 {
		return nil
	}

	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(setMetricQuery)
	if err != nil {
		return err
	}

	for _, m := range meticsSlice {
		if m.MType != metrics.Counter && m.MType != metrics.Gauge {
			return metricerrors.ErrMetricTypeNotImplemented
		}
		_, err := stmt.Exec(m.ID, m.MType, m.Delta, m.Value)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (ps *pgstorage) Get(name string) (val metrics.Metric, ok error) {
	var metric metrics.Metric
	row := ps.db.QueryRow(getMetricQuery, name)
	if err := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value); err != nil {
		return metric, err
	}
	return metric, nil
}

func (ps *pgstorage) GetAll() ([]metrics.Metric, error) {
	rows, err := ps.db.Query(getAllMetricsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allMetrics := []metrics.Metric{}
	for rows.Next() {
		var metric metrics.Metric
		if err := rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value); err != nil {
			return nil, err
		}
		allMetrics = append(allMetrics, metric)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return allMetrics, nil
}

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

func NeedToRetry() func(err error) bool {
	return func(err error) bool {
		var pgErr *pgconn.PgError

		return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ConnectionException
	}
}

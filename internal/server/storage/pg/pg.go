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
		return &newStorage, err
	}
	newStorage.db = db
	return &newStorage, nil
}

func (ps *pgstorage) Setup() error {
	_, err := ps.db.Exec(createTableQuery)
	return err
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

	if metric.MType == metrics.Counter {
		cur, ok := ps.Get(metric.ID)
		if ok {
			newDelta := (*metric.Delta + *cur.Delta)
			metric.Delta = &newDelta
		}
	}

	_, err := ps.db.Exec(setMetricQuery, metric.ID, metric.MType, metric.Delta, metric.Value)
	if err != nil {
		return metric, err
	}
	return metric, nil
}

func (ps *pgstorage) SetAll(meticsSlice []metrics.Metric) error {
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

		if m.MType == metrics.Counter {
			if cur, ok := ps.Get(m.ID); ok {
				newDelta := (*m.Delta + *cur.Delta)
				m.Delta = &newDelta
			}
		}

		_, err := stmt.Exec(m.ID, m.MType, m.Delta, m.Value)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (ps *pgstorage) Get(name string) (val metrics.Metric, ok bool) {
	var metric metrics.Metric
	row := ps.db.QueryRow(getMetricQuery, name)
	if err := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value); err != nil {
		return metric, false
	}
	return metric, true
}

func (ps *pgstorage) GetAll() ([]metrics.Metric, error) {
	allMetrics := make([]metrics.Metric, 0)
	rows, err := ps.db.Query(getAllMetricsQuery)
	if err != nil {
		return allMetrics, err
	}
	defer rows.Close()

	for rows.Next() {
		var metric metrics.Metric
		if err := rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value); err != nil {
			return allMetrics, err
		}
		allMetrics = append(allMetrics, metric)
	}
	if err := rows.Err(); err != nil {
		return allMetrics, err
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

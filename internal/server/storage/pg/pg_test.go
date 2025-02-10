package pg

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
	"gotest.tools/v3/assert"
)

var contextTimeout = 40 * time.Second

func TestStorage_SetMetric(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	storage, err := newPostgresStorage(ctx)
	require.NoError(t, err)

	metric := metrics.NewGaugeMetric("gauge_test", 12.8)

	dbMetric, err := storage.Set(ctx, metric)
	require.NoError(t, err)
	assert.Equal(t, metric.ID, dbMetric.ID)
	assert.Equal(t, *metric.Value, *dbMetric.Value)
	assert.Equal(t, metric.MType, dbMetric.MType)

	dbMetric, err = storage.Get(ctx, metric.ID)
	require.NoError(t, err)
	assert.Equal(t, metric.ID, dbMetric.ID)
	assert.Equal(t, *metric.Value, *dbMetric.Value)
	assert.Equal(t, metric.MType, dbMetric.MType)
}

func newPostgresStorage(ctx context.Context) (*pgstorage, error) {
	dbName := "gophermart"
	dbUser := "user"
	dbPassword := "password"

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	storage, err := NewStorage(ctx, db)
	if err != nil {
		return nil, err
	}

	return storage, err
}

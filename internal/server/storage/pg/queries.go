package pg

const (
	getMetricQuery = `SELECT
			id,
			type,
			delta,
			value
		FROM
			metrics
		WHERE id=$1`

	setMetricQuery = `INSERT INTO metrics
			(id, type, delta, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			id=$1, type=$2, delta=metrics.delta+$3, value=$4`

	getAllMetricsQuery = `SELECT
			id,
			type,
			delta,
			value
		FROM
			metrics`

	createTableQuery = `CREATE TABLE IF NOT EXISTS metrics
		(
			id    varchar(255) PRIMARY KEY, 
			type  varchar(30) NOT NULL, 
			delta bigint, 
			value double precision
		);`
)

// Package sqlmetrics provides Prometheus metrics for database/sql.
package sqlmetrics

import (
	"database/sql"

	prom "github.com/prometheus/client_golang/prometheus"
)

// Collector is a prometheus.Collector that exposes *sql.DB metrics.
type Collector struct {
	db *sql.DB

	maxOpenConnections *prom.Desc

	openConnections *prom.Desc
	inUse           *prom.Desc
	idle            *prom.Desc

	waitCount         *prom.Desc
	waitDuration      *prom.Desc
	maxIdleClosed     *prom.Desc
	maxLifetimeClosed *prom.Desc
}

// NewCollector creates a new collector.
func NewCollector(driver, dbName string, db *sql.DB) *Collector {
	constLabels := prom.Labels{
		"driver": driver,
		"db":     dbName,
	}

	return &Collector{
		db: db,

		maxOpenConnections: prom.NewDesc(
			prom.BuildFQName("go_sql", "connections", "max_open"),
			"Maximum number of open connections to the database.",
			nil, constLabels,
		),

		openConnections: prom.NewDesc(
			prom.BuildFQName("go_sql", "connections", "open"),
			"The number of established connections both in use and idle.",
			nil, constLabels,
		),
		inUse: prom.NewDesc(
			prom.BuildFQName("go_sql", "connections", "in_use"),
			"The number of connections currently in use.",
			nil, constLabels,
		),
		idle: prom.NewDesc(
			prom.BuildFQName("go_sql", "connections", "idle"),
			"The number of idle connections.",
			nil, constLabels,
		),

		waitCount: prom.NewDesc(
			prom.BuildFQName("go_sql", "connections", "wait_count"),
			"The total number of connections waited for.",
			nil, constLabels,
		),
		waitDuration: prom.NewDesc(
			prom.BuildFQName("go_sql", "connections", "wait_duration_seconds"),
			"The total time blocked waiting for a new connection.",
			nil, constLabels,
		),
		maxIdleClosed: prom.NewDesc(
			prom.BuildFQName("go_sql", "connections", "max_idle_closed"),
			"The total number of connections closed due to SetMaxIdleConns.",
			nil, constLabels,
		),
		maxLifetimeClosed: prom.NewDesc(
			prom.BuildFQName("go_sql", "connections", "max_lifetime_closed"),
			"The total number of connections closed due to SetConnMaxLifetime.",
			nil, constLabels,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prom.Desc) {
	ch <- c.maxOpenConnections

	ch <- c.openConnections
	ch <- c.inUse
	ch <- c.idle

	ch <- c.waitCount
	ch <- c.waitDuration
	ch <- c.maxIdleClosed
	ch <- c.maxLifetimeClosed
}

func (c *Collector) Collect(ch chan<- prom.Metric) {
	stats := c.db.Stats()

	ch <- prom.MustNewConstMetric(
		c.maxOpenConnections,
		prom.GaugeValue,
		float64(stats.MaxOpenConnections),
	)

	ch <- prom.MustNewConstMetric(
		c.openConnections,
		prom.GaugeValue,
		float64(stats.OpenConnections),
	)
	ch <- prom.MustNewConstMetric(
		c.inUse,
		prom.GaugeValue,
		float64(stats.InUse),
	)
	ch <- prom.MustNewConstMetric(
		c.idle,
		prom.GaugeValue,
		float64(stats.Idle),
	)

	ch <- prom.MustNewConstMetric(
		c.waitCount,
		prom.CounterValue,
		float64(stats.WaitCount),
	)
	ch <- prom.MustNewConstMetric(
		c.waitDuration,
		prom.CounterValue,
		stats.WaitDuration.Seconds(),
	)
	ch <- prom.MustNewConstMetric(
		c.maxIdleClosed,
		prom.CounterValue,
		float64(stats.MaxIdleClosed),
	)
	ch <- prom.MustNewConstMetric(
		c.maxLifetimeClosed,
		prom.CounterValue,
		float64(stats.MaxLifetimeClosed),
	)
}

// check interfaces
var (
	_ prom.Collector = (*Collector)(nil)
)

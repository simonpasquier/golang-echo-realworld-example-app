package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type StoreMetrics struct {
	total    *prometheus.CounterVec
	failed   *prometheus.CounterVec
	duration prometheus.ObserverVec
}

const (
	opCreate = "create"
	opDelete = "delete"
	opRead   = "read"
	opUpdate = "update"
)

func NewStoreMetrics(r prometheus.Registerer) *StoreMetrics {
	total := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "store_operations_total",
			Help: "Total number of store operations.",
		},
		[]string{"operation"},
	)
	failed := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "store_operations_failed_total",
			Help: "Total number of store operations that failed.",
		},
		[]string{"operation"},
	)
	duration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "store_operation_seconds",
			Help:    "A histogram of latencies for store requests.",
			Buckets: []float64{.05, .1, .25, .5, .75, 1, 2, 5},
		},
		[]string{"operation"},
	)
	r.MustRegister(total, failed, duration)

	for _, o := range []string{opCreate, opDelete, opRead, opUpdate} {
		total.WithLabelValues(o)
		failed.WithLabelValues(o)
		duration.WithLabelValues(o)
	}

	return &StoreMetrics{
		total:    total,
		failed:   failed,
		duration: duration,
	}
}
func (m *StoreMetrics) wrapRequest(op string, f func() bool) {
	now := time.Now()
	m.total.WithLabelValues(op).Inc()
	if !f() {
		m.failed.WithLabelValues(op).Inc()
	}
	m.duration.WithLabelValues(op).Observe(time.Since(now).Seconds())
}

package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/xesina/golang-echo-realworld-example-app/db"
	"github.com/xesina/golang-echo-realworld-example-app/handler"
	"github.com/xesina/golang-echo-realworld-example-app/metrics"
	"github.com/xesina/golang-echo-realworld-example-app/router"
	"github.com/xesina/golang-echo-realworld-example-app/store"
)

const appName = "realworld"

// These variables are injected at build time.
var appVersion, appRevision, appBranch string

func main() {
	buildInfo := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: appName,
			Name:      "build_info",
			Help:      "A metric with a constant '1' value labeled by version, branch and revision",
			ConstLabels: prometheus.Labels{
				"branch":   appBranch,
				"revision": appRevision,
				"version":  appVersion,
			},
		},
	)
	buildInfo.Set(1)
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
		buildInfo,
	)
	r := router.New(reg, appName)
	v1 := r.Group("/api")

	d := db.New()
	db.AutoMigrate(d)

	m := metrics.NewStoreMetrics(reg, appName)
	us := metrics.NewUserStore(store.NewUserStore(d), m)
	as := metrics.NewArticleStore(store.NewArticleStore(d), m)
	h := handler.NewHandler(us, as)
	h.Register(v1)
	r.Logger.Fatal(r.Start("127.0.0.1:8585"))
}

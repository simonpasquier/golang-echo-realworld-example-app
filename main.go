package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/xesina/golang-echo-realworld-example-app/db"
	"github.com/xesina/golang-echo-realworld-example-app/handler"
	"github.com/xesina/golang-echo-realworld-example-app/metrics"
	"github.com/xesina/golang-echo-realworld-example-app/router"
	"github.com/xesina/golang-echo-realworld-example-app/store"
)

// These variables are injected at build time.
var appVersion, appRevision, appBranch string

func main() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	// Prefix all application metrics.
	appReg := prometheus.WrapRegistererWithPrefix("realworld_", reg)
	buildInfo := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "build_info",
			Help: "A metric with a constant '1' value labeled by version, branch and revision",
			ConstLabels: prometheus.Labels{
				"branch":   appBranch,
				"revision": appRevision,
				"version":  appVersion,
			},
		},
	)
	buildInfo.Set(1)
	appReg.MustRegister(buildInfo)

	r := router.New(reg)
	v1 := r.Group("/api")

	d := db.New()
	db.AutoMigrate(d)

	m := metrics.NewStoreMetrics(appReg)
	us := metrics.NewUserStore(store.NewUserStore(d), m)
	as := metrics.NewArticleStore(store.NewArticleStore(d), m)
	h := handler.NewHandler(us, as)
	h.Register(v1)
	r.Logger.Fatal(r.Start("127.0.0.1:8585"))
}

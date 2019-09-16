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

func main() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
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

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/xesina/golang-echo-realworld-example-app/db"
	"github.com/xesina/golang-echo-realworld-example-app/handler"
	"github.com/xesina/golang-echo-realworld-example-app/metrics"
	"github.com/xesina/golang-echo-realworld-example-app/router"
	"github.com/xesina/golang-echo-realworld-example-app/store"
)

// These variables are injected at build time.
var appVersion, appRevision, appBranch string

type faultMiddleware struct {
	errRatio   float64
	delayRatio float64
	delay      float64

	mtx  sync.Mutex
	rand *rand.Rand
}

func newFaultMiddleware(errRatio, delayRatio float64, delay time.Duration) *faultMiddleware {
	return &faultMiddleware{
		errRatio:   errRatio,
		delayRatio: delayRatio,
		delay:      float64(delay),
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (f *faultMiddleware) gotError() bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return f.rand.Float64() < f.errRatio
}

func (f *faultMiddleware) gotDelay() bool {
	if f.delay <= 0 {
		return false
	}
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return f.rand.Float64() < f.delayRatio
}

func (f *faultMiddleware) addDelay() {
	if !f.gotDelay() {
		return
	}
	f.mtx.Lock()
	r := float64(f.rand.Int63n(int64(f.delay * 0.2)))
	f.mtx.Unlock()
	time.Sleep(time.Duration(f.delay*0.9 + r))
}

func (f *faultMiddleware) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if f.gotError() {
			return errors.New("got error (fault middleware)")
		}
		f.addDelay()
		return next(c)
	}
}

func (f *faultMiddleware) String() string {
	return fmt.Sprintf("mean error ratio=%f, mean delay ratio= %f, mean delay duration=%v", f.errRatio, f.delayRatio, time.Duration(f.delay))
}

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
	var (
		errRatio   float64
		delayRatio float64 = 1
		delay      time.Duration
	)
	if s := os.Getenv("ERROR_RATIO"); s != "" {
		var err error
		errRatio, err = strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
	}
	if s := os.Getenv("DELAY_RATIO"); s != "" {
		var err error
		delayRatio, err = strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
	}
	if s := os.Getenv("DELAY"); s != "" {
		var err error
		delay, err = time.ParseDuration(s)
		if err != nil {
			panic(err)
		}
	}
	r.Logger.Infoj(map[string]interface{}{
		"msg":          "fault inhection settings",
		"error_ratio":  errRatio,
		"delay_ration": delayRatio,
		"delay":        delay,
	})
	fault := newFaultMiddleware(errRatio, delayRatio, delay)
	v1 := r.Group("/api", fault.Process)

	dbDriver, dbHost, dbPort, dbUser, dbPassword, dbName := os.Getenv("DB_DRIVER"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")
	if dbDriver == "" {
		dbDriver = "sqlite3"
	}
	r.Logger.Infoj(map[string]interface{}{
		"msg":      "database settings",
		"driver":   dbDriver,
		"host":     dbHost,
		"port":     dbPort,
		"name":     dbName,
		"user":     dbUser,
		"password": "****",
	})
	var datasource string
	switch dbDriver {
	case "postgres":
		if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
			panic("need to set DB_HOST, DB_PORT, DB_USER, DB_PASSWORD and DB_NAME environment variables")
		}
		datasource = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", dbHost, dbPort, dbUser, dbPassword, dbName)
	case "sqlite3":
		if dbName == "" {
			dbName = ":memory:"
		}
		datasource = dbName
		datasource = "./realworld.db"
	default:
		panic(fmt.Sprintf("unsupported database driver: %s", dbDriver))
	}

	d := db.New(dbDriver, datasource)
	db.AutoMigrate(d)

	m := metrics.NewStoreMetrics(appReg)
	us := metrics.NewUserStore(store.NewUserStore(d), m)
	as := metrics.NewArticleStore(store.NewArticleStore(d), m)
	h := handler.NewHandler(us, as)
	h.Register(v1)
	r.Logger.Fatal(r.Start(":8585"))
}

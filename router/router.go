package router

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(reg *prometheus.Registry) *echo.Echo {
	duration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "A histogram of latencies for HTTP requests.",
			Buckets: []float64{.1, .2, .3, .4, .5, .7, .9, 1, 1.5, 3},
		},
		[]string{"code", "method"},
	)
	reg.MustRegister(duration)

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			now := time.Now()
			err := next(c)
			// Don't track duration metric for /metrics.
			if c.Path() != "/metrics" {
				duration.WithLabelValues(strconv.Itoa(c.Response().Status), c.Request().Method).Observe(time.Since(now).Seconds())
			}
			return err
		}
	})
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.GET("/metrics",
		func(c echo.Context) error {
			promhttp.InstrumentMetricHandler(
				reg,
				promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
			).ServeHTTP(c.Response(), c.Request())
			return nil
		},
	)
	e.Validator = NewValidator()
	return e
}

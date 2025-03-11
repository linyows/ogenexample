package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/XSAM/otelsql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/linyows/ogenexample/db/dbgen"
	"github.com/linyows/ogenexample/oas/oasgen"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var _ oasgen.Handler = (*handler)(nil)

type handler struct {
	q  *dbgen.Queries
	tp *trace.TracerProvider
}

func Server() (*oasgen.Server, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	tp := setupDebugTracerProvider()
	mp := setupDebugMeterProvider()
	// Set to global
	//otel.SetTracerProvider(tp)
	//otel.SetMeterProvider(mp)

	return oasgen.NewServer(
		&handler{
			q:  dbgen.New(db),
			tp: tp,
		},
		customNotFound(),
		customErrorHandler(),
		oasgen.WithTracerProvider(tp),
		oasgen.WithMeterProvider(mp),
	)
}

func connectDB() (*sql.DB, error) {
	dsn := "root@tcp(localhost:3306)/ogenexample?parseTime=true"
	db, err := otelsql.Open("mysql", dsn, otelsql.WithAttributes(semconv.DBSystemMySQL))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)
	db.SetConnMaxIdleTime(0)

	if err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(semconv.DBSystemMySQL)); err != nil {
		return nil, err
	}

	return db, nil
}

func customNotFound() oasgen.ServerOption {
	return oasgen.WithNotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"error": "Not found!!!"}`)
	})
}

func customErrorHandler() oasgen.ServerOption {
	return oasgen.WithErrorHandler(func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
		var (
			code    = http.StatusInternalServerError
			ogenErr ogenerrors.Error
		)
		switch {
		case errors.Is(err, ht.ErrNotImplemented):
			code = http.StatusNotImplemented
		case errors.As(err, &ogenErr):
			code = ogenErr.Code()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_, _ = io.WriteString(w, fmt.Sprintf(`{"error": "%s"}`, http.StatusText(code)))
	})
}

const (
	appName = "ogenexample"
	appVer  = "v1.2.3"
	appEnv  = "dev"
)

func setupDebugTracerProvider() *trace.TracerProvider {
	//exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	exp, err := stdouttrace.New()
	if err != nil {
		log.Fatalf("Failed to create stdout exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(appName),
			semconv.ServiceVersion(appVer),
			semconv.DeploymentEnvironment(appEnv),
		)),
	)
	return tp
}

func setupDebugMeterProvider() *metric.MeterProvider {
	//exporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	exporter, err := stdoutmetric.New()
	if err != nil {
		log.Fatalf("Failed to create stdout metric exporter: %v", err)
	}

	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(appName),
			semconv.ServiceVersion(appVer),
			semconv.DeploymentEnvironment(appEnv),
		))
	if err != nil {
		log.Fatalf("Failed to merge resource: %v", err)
	}

	reader := metric.NewPeriodicReader(exporter)
	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(reader),
	)
	return mp
}

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/bzon/errors"
	"github.com/go-kit/kit/log"
	"go.opencensus.io/trace"
)

func main() {
	// logging
	var logger log.Logger
	{
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
		logger = log.With(logger, "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)
	}

	// tracing
	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     "localhost:6831",
		CollectorEndpoint: "http://localhost:14268/api/traces",
		ServiceName:       "erroring-service",
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create the Jaeger exporter: %v", err))
	}
	trace.RegisterExporter(je)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	for {
		workerr := work(context.Background(), logger)
		if ec, ok := workerr.(errors.Error); ok {
			logger.Log(
				"message", ec.Error(),
				"logging.googleapis.com/spanId", ec.TraceContext().SpanID,
				"logging.googleapis.com/trace", ec.TraceContext().TraceID,
				"logging.googleapis.com/sourceLocation", ec.SourceLocation(),
			)
		}
		time.Sleep(3 * time.Second)
	}

}

func work(ctx context.Context, logger log.Logger) error {
	_, span := trace.StartSpan(ctx, "work")
	defer span.End()

	// error with context
	err := errors.NewT(span, "error")
	return err
}

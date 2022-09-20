package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	logs "go.opentelemetry.io/otel/logs2"
	"go.opentelemetry.io/otel/trace"
)

func mainAlt() {
	log.Printf("Waiting for connection...")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer := otel.Tracer("test-tracer")

	// Attributes represent additional key-value descriptors that can be bound
	// to a metric observer or recorder.
	commonAttrs := []attribute.KeyValue{
		attribute.String("attrA", "chocolate"),
		attribute.String("attrB", "raspberry"),
		attribute.String("attrC", "vanilla"),
	}

	// Setup Otel Logger
	logger := logs.New(logs.NewOtelHandler(os.Stdout))
	ctx = logs.NewContext(ctx, logger)

	// work begins
	ctx, span := tracer.Start(
		ctx,
		"CollectorExporter-Example",
		trace.WithAttributes(commonAttrs...),
	)
	defer span.End()

	for i := 0; i < 10; i++ {
		doHardWork2(ctx, tracer, i)
	}

	log.Printf("Done!")
}

func doHardWork2(ctx context.Context, tracer trace.Tracer, i int) {
	logger := logs.FromContext(ctx)
	logger.LogAttrs(ctx, logs.InfoLevel, "Begin doing really hard work", logs.Any("iter", i+1))

	_, span := tracer.Start(ctx, fmt.Sprintf("Sample-%d", i))

	doInnerWork(ctx, i)

	span.End()
}

func doInnerWork2(ctx context.Context, i int) {
	logger := logs.FromContext(ctx)
	logger.LogAttrs(ctx, logs.InfoLevel, "Continue doing really hard work", logs.Any("iter", i+1))
	<-time.After(time.Second)
}

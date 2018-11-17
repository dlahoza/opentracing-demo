package mysql

import (
	"context"
	"math/rand"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Database struct {
	workers chan struct{}
	tracer  opentracing.Tracer
}

func New(workers int, tracer opentracing.Tracer) *Database {
	return &Database{workers: make(chan struct{}, workers), tracer: tracer} //Max 5 workers
}

func (d *Database) Query(ctx context.Context, query string) {
	d.workers <- struct{}{}
	defer func() { <-d.workers }()
	// simulate opentracing instrumentation of an SQL query
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := d.tracer.StartSpan("SQL SELECT", opentracing.ChildOf(span.Context()))
		ext.SpanKindRPCClient.Set(span)
		ext.PeerService.Set(span, "mysql")
		span.SetTag("sql.query", query)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)+50)) // 50-100ms query duration
}

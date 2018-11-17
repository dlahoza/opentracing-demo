package redis

import (
	"context"
	"math/rand"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
)

type Redis struct {
	tracer opentracing.Tracer
}

func New(tracer opentracing.Tracer) *Redis {
	return &Redis{tracer: tracer} //Max 5 workers
}

func (d *Redis) Get(ctx context.Context, key string) {
	// simulate opentracing instrumentation of an SQL query
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := d.tracer.StartSpan("Redis Get", opentracing.ChildOf(span.Context()))
		ext.SpanKindRPCClient.Set(span)
		ext.PeerService.Set(span, "redis")
		span.SetTag("redis.get", key)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)+50)) // 50-100ms query duration
}

func (d *Redis) Set(ctx context.Context, key string, errorProbability int) error {
	// simulate opentracing instrumentation of an SQL query
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		span := d.tracer.StartSpan("Redis Set", opentracing.ChildOf(span.Context()))
		ext.SpanKindRPCClient.Set(span)
		ext.PeerService.Set(span, "redis")
		span.SetTag("redis.set", key)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(400)+100)) // 100-500ms query duration
	if rand.Intn(100) > 100-errorProbability {
		ext.Error.Set(span, true)
		return errors.New("Redis error")
	}
	return nil
}

package tracing

import (
	"fmt"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
)

// HTTPClient wraps an http.Client with tracing instrumentation.
type HTTPClient struct {
	Tracer opentracing.Tracer
	Client *http.Client
}

// GetJSON executes HTTP GET against specified url and tried to parse
// the response into out object.
func (c *HTTPClient) Do(req *http.Request) (response *http.Response, err error) {
	if span := opentracing.SpanFromContext(req.Context()); span != nil {

		// start a new Span to wrap HTTP request
		newspan := c.Tracer.StartSpan(
			fmt.Sprintf("Client HTTP %s: %s", req.Method, req.URL.EscapedPath()),
			opentracing.ChildOf(span.Context()),
		)

		// make sure the Span is finished once we're done
		defer newspan.Finish()
		c.Tracer.Inject(newspan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	}

	return c.Client.Do(req)
}

func TracingMiddleware(tr opentracing.Tracer) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return nethttp.MiddlewareFunc(tr, handler.ServeHTTP, nethttp.OperationNameFunc(func(r *http.Request) string {
			return "HTTP " + r.Method + " " + r.URL.EscapedPath()
		}))
	}
}

package main

import (
	"github.com/opentracing/opentracing-go"
	"log"
	"net/http"

	"github.com/DLag/opentracing-demo/pkg/redis"
	"github.com/DLag/opentracing-demo/pkg/tracing"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()
	tracer := tracing.Init("authority")
	opentracing.SetGlobalTracer(tracer)
	s := AuthorityService{cache: redis.New(tracer)}
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(tracing.TracingMiddleware(tracer))

	r.Post("/auth", s.AuthHandler)
	r.Post("/check", s.CheckHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("{\"status\": \"OK\"}"))
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

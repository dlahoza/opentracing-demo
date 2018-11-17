package main

import (
	"github.com/opentracing/opentracing-go"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/DLag/opentracing-demo/pkg/mysql"
	"github.com/DLag/opentracing-demo/pkg/tracing"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()
	tracer := tracing.Init("treasury")
	opentracing.SetGlobalTracer(tracer)

	rand.Seed(time.Now().UnixNano())

	s := AuthorityService{db: mysql.New(10, tracer)}
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(tracing.TracingMiddleware(tracer))

	r.Post("/balance", s.BalanceHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("{\"status\": \"OK\"}"))
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

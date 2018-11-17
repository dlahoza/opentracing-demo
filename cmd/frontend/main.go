package main

import (
	"log"
	"net/http"

	"github.com/DLag/opentracing-demo/pkg/tracing"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	tracer := tracing.Init("frontend")

	client := &tracing.HTTPClient{
		Tracer: tracer,
		Client: http.DefaultClient,
	}
	s := &FrontendService{c: client}
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(tracing.TracingMiddleware(tracer))

	r.Get("/auth", s.AuthHandler)
	r.Get("/balance", s.BalanceHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("{\"status\": \"OK\"}"))
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

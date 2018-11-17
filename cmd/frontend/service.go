package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DLag/opentracing-demo/pkg/models"
	"github.com/DLag/opentracing-demo/pkg/tracing"
	"github.com/opentracing/opentracing-go"
)

type FrontendService struct {
	c *tracing.HTTPClient
}

func (s *FrontendService) AuthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.FormValue("user")

	if userID == "" {
		http.Error(w, "Missing required 'user' parameter", http.StatusBadRequest)
		return
	}

	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.SetTag("userID", userID)
	}

	data, _ := json.Marshal(models.Check{User: userID})

	req, _ := http.NewRequest("POST", "http://authority:8080/auth", bytes.NewBuffer(data))

	resp, err := s.c.Do(req.WithContext(ctx))

	if err != nil || (resp.StatusCode >= 400 && resp.StatusCode < 500) {
		http.Error(w, "Wrong parameters", http.StatusBadRequest)
		return
	}

	if resp.StatusCode >= 500 {
		http.Error(w, "Error on subrequest", http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusInternalServerError)
		return
	}

	check := new(models.Check)
	err = json.Unmarshal(buf, check)
	if err != nil || check.Session == "" {
		http.Error(w, "Cannot unmarshal response of subrequest", http.StatusInternalServerError)
		return
	}

	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.SetTag("sessionID", check.Session)
	}

	w.Write([]byte(check.Session))
}

func (s *FrontendService) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := r.FormValue("session")

	if sessionID == "" {
		http.Error(w, "Missing required 'session' parameter", http.StatusBadRequest)
		return
	}

	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.SetTag("sessionID", sessionID)
	}

	data, _ := json.Marshal(models.Check{Session: sessionID})

	req, _ := http.NewRequest("POST", "http://authority:8080/check", bytes.NewBuffer(data))

	resp, err := s.c.Do(req.WithContext(ctx))

	if err != nil || (resp.StatusCode >= 400 && resp.StatusCode < 500) {
		http.Error(w, "Wrong parameters", http.StatusBadRequest)
		return
	}

	if resp.StatusCode >= 500 {
		http.Error(w, "Error on subrequest", http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusInternalServerError)
		return
	}

	check := new(models.Check)
	err = json.Unmarshal(buf, check)
	if err != nil || check.User == "" {
		http.Error(w, "Cannot unmarshal response of subrequest", http.StatusInternalServerError)
		return
	}

	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.SetTag("userID", check.User)
	}

	userID := check.User

	// TREASURY

	data, _ = json.Marshal(models.Check{User: userID})

	req, _ = http.NewRequest("POST", "http://treasury:8080/balance", bytes.NewBuffer(data))

	resp, err = s.c.Do(req.WithContext(ctx))

	if err != nil || (resp.StatusCode >= 400 && resp.StatusCode < 500) {
		http.Error(w, "Wrong parameters", http.StatusBadRequest)
		return
	}

	if resp.StatusCode >= 500 {
		http.Error(w, "Error on subrequest", http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()
	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusInternalServerError)
		return
	}

	balance := new(models.Balance)
	err = json.Unmarshal(buf, balance)
	if err != nil || balance.Balance == 0 {
		http.Error(w, "Cannot unmarshal response of subrequest", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, balance.Balance)
}

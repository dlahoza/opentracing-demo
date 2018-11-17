package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/DLag/opentracing-demo/pkg/models"
	"github.com/DLag/opentracing-demo/pkg/redis"

	"github.com/google/uuid"
)

type AuthorityService struct {
	cache    *redis.Redis
	sessions sync.Map
}

func (s *AuthorityService) AuthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusInternalServerError)
		return
	}

	check := new(models.Check)
	err = json.Unmarshal(buf, check)
	if err != nil || check.User == "" {
		http.Error(w, "Missing required 'session' parameter", http.StatusBadRequest)
		return
	}

	u := uuid.New()
	sessionID := "session-" + u.String()

	if err := s.cache.Set(ctx, sessionID, 10); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.sessions.Store(sessionID, check.User)

	data, _ := json.Marshal(models.Check{Session: sessionID})
	w.Write(data)
}

func (s *AuthorityService) CheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusInternalServerError)
		return
	}

	check := new(models.Check)
	err = json.Unmarshal(buf, check)
	if err != nil || check.Session == "" {
		http.Error(w, "Missing required 'session' parameter", http.StatusBadRequest)
		return
	}

	s.cache.Get(ctx, check.Session)

	if userID, ok := s.sessions.Load(check.Session); ok {
		data, _ := json.Marshal(models.Check{Session: check.Session, User: userID.(string)})
		w.Write(data)
		return
	}

	w.WriteHeader(http.StatusForbidden)

}

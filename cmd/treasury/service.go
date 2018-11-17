package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"

	"github.com/DLag/opentracing-demo/pkg/models"
	"github.com/DLag/opentracing-demo/pkg/mysql"
)

type AuthorityService struct {
	db       *mysql.Database
	sessions sync.Map
}

func (s *AuthorityService) BalanceHandler(w http.ResponseWriter, r *http.Request) {
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

	balance := rand.Intn(1000000000000)

	s.db.Query(ctx, fmt.Sprintf("SELECT from money where user = %q", check.User))

	data, _ := json.Marshal(models.Balance{Balance: balance})
	w.Write(data)
}

package community

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"net/http"
)

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger, secret string) {
	res := resource{service, logger}
	r.HandleFunc("/api/v1/communities", res.createCommunityHandler).Methods("POST")

	r.Use()
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) createCommunityHandler(w http.ResponseWriter, req *http.Request) {
	var input CreateCommunityRequest

	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("calling the service layer")
	community, err := r.service.Create(req.Context(), input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(community)
}

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
	r.HandleFunc("/api/v1/communities", res.GetAllCommunitiesHandler).Methods("GET")
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

func (r resource) GetAllCommunitiesHandler(w http.ResponseWriter, req *http.Request) {
	results, err := r.service.GetAllCommunities(req.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

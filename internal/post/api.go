package post

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger, secret string) {
	res := resource{service, logger}
	//r.HandleFunc("/api/v1/categories/{id}", res.getByIdHandler).Methods("GET")
	r.HandleFunc("/api/v1/posts", res.getAllHandler).Methods("GET")

	// Protected Endpoints
	//r.Handle("/api/v1/categories", auth.AuthenticateMiddleware(auth.RoleMiddleware(http.HandlerFunc(res.create), "admin"), secret)).Methods("POST")
	r.Use()
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) getAllHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	json.NewEncoder(w).Encode(id)
}

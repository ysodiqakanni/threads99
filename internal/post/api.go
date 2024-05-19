package post

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger, secret string) {
	res := resource{service, logger}
	//r.HandleFunc("/api/v1/categories/{id}", res.getByIdHandler).Methods("GET")
	//r.HandleFunc("/api/v1/posts", res.getAllHandler).Methods("GET")
	r.HandleFunc("/api/v1/posts", res.getByIdHandler).Methods("GET")

	// Protected Endpoints
	//r.Handle("/api/v1/categories", auth.AuthenticateMiddleware(auth.RoleMiddleware(http.HandlerFunc(res.create), "admin"), secret)).Methods("POST")
	r.Use()
}

type resource struct {
	service Service
	logger  log.Logger
}

//func NewService(repo Repository, userRepo user.Repository, logger log.Logger) Service {
//	return service{repo, userRepo, logger}
//}

func (r resource) getByIdHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	idk, _ := primitive.ObjectIDFromHex(id)

	post, _ := r.service.Get(req.Context(), idk)
	json.NewEncoder(w).Encode(post)
}

func (r resource) getAllHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	json.NewEncoder(w).Encode(id)
}
func (r resource) createNewPostHandler(w http.ResponseWriter, req *http.Request) {
	var input CreateNewPostRequest
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

	r.service
}

func (r resource) createCommentHandler() {

}
func (r resource) upVotePostHandler() {

}
func (r resource) downVotePostHandler() {

}
func (r resource) upVoteCommentHandler() {

}
func (r resource) downVoteCommentHandler() {

}

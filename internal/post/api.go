package post

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ysodiqakanni/threads99/internal/models"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger, secret string) {
	res := resource{service, logger}
	//r.HandleFunc("/api/v1/categories/{id}", res.getByIdHandler).Methods("GET")
	//r.HandleFunc("/api/v1/posts", res.getAllHandler).Methods("GET")
	//r.HandleFunc("/api/v1/posts", res.getByIdHandler).Methods("GET")
	r.HandleFunc("/api/v1/posts/recent", res.GetAllRecentPostsHandler).Methods("GET")
	r.HandleFunc("/api/v1/posts/{postId}", res.getByIdHandler).Methods("GET")

	r.HandleFunc("/api/v1/posts/{postId}/comments", res.getCommentsHandler).Methods("GET")
	r.HandleFunc("/api/v1/posts", res.createNewPostHandler).Methods("POST")
	r.HandleFunc("/api/v1/posts/add-comment", res.createCommentHandler).Methods("PUT")
	r.HandleFunc("/api/v1/posts/upvote-comment", res.voteCommentHandler).Methods("PUT")
	r.HandleFunc("/api/v1/posts/vote", res.votePostHandler).Methods("PUT")
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

func (r resource) getByIdHandlerOld(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	idk, _ := primitive.ObjectIDFromHex(id)

	post, _ := r.service.Get(req.Context(), idk)
	json.NewEncoder(w).Encode(post)
}

func (r resource) getByIdHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["postId"]
	idk, _ := primitive.ObjectIDFromHex(id)

	post, _ := r.service.GetPostLiteById(req.Context(), idk)

	response := models.NewSuccessResponse(
		post,
		"Post retrieved successfully",
	)
	json.NewEncoder(w).Encode(response)
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
		response := models.NewErrorResponse(
			[]string{err.Error()},
			"Invalid input",
		)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return

		//http.Error(w, err.Error(), http.StatusBadRequest)
		//return
	}
	if err := input.Validate(); err != nil {
		r.logger.With(req.Context()).Info(err)
		response := models.NewErrorResponse(
			[]string{err.Error()},
			"Validation failed!",
		)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
		//http.Error(w, err.Error(), http.StatusBadRequest)
		//return
	}

	err, postId := r.service.CreatePost(req.Context(), input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		response := models.NewErrorResponse(
			[]string{err.Error()},
			"An unknown error occurred!",
		)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
		//
		//http.Error(w, "Error creating new post "+err.Error(), http.StatusInternalServerError)
		//return
	}

	// Todo: should this endpoint return the new post ID on success?
	response := models.NewSuccessResponse(
		postId,
		"Post created!",
	)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (r resource) GetAllRecentPostsHandler(w http.ResponseWriter, req *http.Request) {
	results, err := r.service.GetAllRecentPosts(req.Context())
	if err != nil {
		response := models.NewErrorResponse(
			[]string{err.Error()},
			"Failed to fetch posts",
		)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.NewSuccessResponse(
		results,
		"Posts retrieved successfully",
	)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (r resource) createCommentHandler(w http.ResponseWriter, req *http.Request) {
	var input AddCommentToPostRequest
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

	err = r.service.AddCommentToPost(req.Context(), input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, "Error creating new comment "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (r resource) votePostHandler(w http.ResponseWriter, req *http.Request) {
	var input PostUpvoteRequest
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = input.Validate()
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.service.UpvotePost(req.Context(), input)
}

func (r resource) voteCommentHandler(w http.ResponseWriter, req *http.Request) {
	var input CommentUpvoteRequest
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = input.Validate()
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.service.UpvoteComment(req.Context(), input)
}

func (r resource) getCommentsHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	postId := vars["postId"]
	fmt.Println("Attempting to load comments for post: " + postId)
	results, err := r.service.GetCommentsByPostId(req.Context(), postId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

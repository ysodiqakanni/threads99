package comment

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ysodiqakanni/threads99/internal/auth"
	"github.com/ysodiqakanni/threads99/internal/dto"
	"github.com/ysodiqakanni/threads99/internal/helper"
	"github.com/ysodiqakanni/threads99/internal/models"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger, secret string) {
	res := resource{service, logger}
	r.HandleFunc("/api/v1/posts/{postId}/comments", res.getCommentsByPostIdHandler).Methods("GET")

	// Protected Endpoints
	r.Handle("/api/v1/comments", auth.AuthenticateMiddleware(http.HandlerFunc(res.createCommentHandler),
		secret)).Methods("POST")
	r.Use()
}

type resource struct {
	service Service
	logger  log.Logger
}

//func NewService(repo Repository, userRepo user.Repository, logger log.Logger) Service {
//	return service{repo, userRepo, logger}
//}

func (r resource) createCommentHandler(w http.ResponseWriter, req *http.Request) {
	var input dto.CreateNewCommentRequest
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId, ok := req.Context().Value("userId").(string)
	username, ok1 := req.Context().Value("username").(string)
	if !ok || !ok1 {
		// Handle case where userId is not in context
		helper.EncodeErrorResponse(w, errors.New("Session Expired."),
			"", "401")
		return
	}
	input.CreatedByUserId = userId
	input.CreatedByUserName = username
	if err := input.Validate(); err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err, commentId := r.service.CreateNewComment(req.Context(), input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		response := models.NewErrorResponse(
			[]string{err.Error()},
			"An unknown error occurred while creating comment.!"+err.Error(), "500",
		)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.NewSuccessResponse(
		commentId,
		"Comment created!",
	)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (r resource) getCommentsByPostIdHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	postId := vars["postId"]
	fmt.Println("Attempting to load comments for post: " + postId)

	_id, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		response := models.NewErrorResponse(
			[]string{err.Error()},
			"Failed to fetch comments. Invalid postId", "400",
		)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	commentsResult, err := r.service.GetCommentsByPostId(req.Context(), _id)

	if err != nil {
		response := models.NewErrorResponse(
			[]string{err.Error()},
			"Error fetching comments.", "500",
		)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.NewSuccessResponse(
		commentsResult,
		"Comments retrieved successfully",
	)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return
}

/*
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
*/

/*
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
*/

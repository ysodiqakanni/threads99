package community

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ysodiqakanni/threads99/internal/auth"
	"github.com/ysodiqakanni/threads99/internal/helper"
	"github.com/ysodiqakanni/threads99/internal/models"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger, secret string) {
	res := resource{service, logger}
	//r.HandleFunc("/api/v1/communities", res.createCommunityHandler).Methods("POST")

	r.HandleFunc("/api/v1/communities", res.GetAllCommunitiesHandler).Methods("GET")
	r.HandleFunc("/api/v1/communities/{id}", res.getByIdHandler).Methods("GET")

	// Protected endpoints.
	r.Handle("/api/v1/communities", auth.AuthenticateMiddleware(http.HandlerFunc(res.createCommunityHandler), secret)).Methods("POST")

	r.Use()
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) getByIdHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	idk, _ := primitive.ObjectIDFromHex(id)

	community, _ := r.service.Get(req.Context(), idk)

	response := models.NewSuccessResponse(
		community,
		"Community retrieved successfully",
	)
	json.NewEncoder(w).Encode(response)
}
func (r resource) createCommunityHandler(w http.ResponseWriter, req *http.Request) {
	var input CreateCommunityRequest

	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		helper.EncodeErrorResponse(w, errors.New("Invalid Model"),
			"", "400")
		return
	}
	if err := input.Validate(); err != nil {
		r.logger.With(req.Context()).Info(err)
		helper.EncodeErrorResponse(w, errors.New("Invalid Model"),
			"", "400")
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
	fmt.Println(username)
	fmt.Println("calling the service layer")
	input.CreatedByUserId = userId
	community, err := r.service.Create(req.Context(), input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		helper.EncodeErrorResponse(w, err,
			"", "500")
		return
	}

	helper.EncodeSuccessResponse(w, community, "Community created successfully.")
}

func (r resource) GetAllCommunitiesHandler(w http.ResponseWriter, req *http.Request) {
	results, err := r.service.GetAllCommunities(req.Context())
	if err != nil {
		helper.EncodeApiFailureResponse(w, err, "Failed to fetch communities", "500")
		return
	}

	response := models.NewSuccessResponse(
		results,
		"Communities retrieved successfully",
	)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

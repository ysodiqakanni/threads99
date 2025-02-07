package user

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ysodiqakanni/threads99/internal/dto"
	"github.com/ysodiqakanni/threads99/internal/entity"
	"github.com/ysodiqakanni/threads99/internal/models"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger) {
	res := resource{service, logger}
	r.HandleFunc("/api/v1/users/register", res.registerUserHandler).Methods("POST")

}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) registerUserHandler(w http.ResponseWriter, req *http.Request) {
	var input dto.CreateNewUserRequestDto
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		response := models.NewErrorResponse(
			[]string{err.Error()},
			"Bad data!",
		)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
	if err := input.Validate(); err != nil {
		r.logger.With(req.Context()).Info(err)
		response := models.NewErrorResponse(
			[]string{err.Error()},
			"Bad Request",
		)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	userObj, err := r.service.Create(req.Context(), input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		response := models.NewErrorResponse(
			[]string{err.Error()},
			"User creation failed.",
		)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	userIdentity := entity.User{
		Email:    input.Email,
		ID:       userObj.UserObjectId,
		Username: input.Email,
	}

	token, err := r.service.GenerateJWT(userIdentity)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		response := models.NewErrorResponse(
			[]string{"User created but error logging in. Try to refresh and try again "},
			"Internal Server error",
		)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	userObj.Token = token
	response := models.NewSuccessResponse(
		userObj,
		"User registered.",
	)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (r resource) getByIdHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	idk, _ := primitive.ObjectIDFromHex(id)

	category, _ := r.service.Get(req.Context(), idk)
	json.NewEncoder(w).Encode(category)
}

func (r resource) createNewUser(w http.ResponseWriter, req *http.Request) {
	// check if user already exists

}

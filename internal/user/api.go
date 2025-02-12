package user

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ysodiqakanni/threads99/internal/dto"
	"github.com/ysodiqakanni/threads99/internal/entity"
	"github.com/ysodiqakanni/threads99/internal/helper"
	"github.com/ysodiqakanni/threads99/internal/models"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger) {
	res := resource{service, logger}
	r.HandleFunc("/api/v1/auth/login", res.loginHandler).Methods("POST")
	r.HandleFunc("/api/v1/user/register", res.registerUserHandler).Methods("POST")
	r.HandleFunc("/api/v1/user/profile/{userId}", res.getUserPublicMetadata).Methods("GET")
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
		helper.EncodeErrorResponse(w, err, err.Error(), "400")
		return
	}
	if err := input.Validate(); err != nil {
		r.logger.With(req.Context()).Info(err)
		helper.EncodeErrorResponse(w, err, err.Error(), "400")
		return
	}

	userObj, err := r.service.Create(req.Context(), input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		helper.EncodeErrorResponse(w, err, err.Error(), "500")
		return
	}

	userIdentity := entity.User{
		Email:    input.Email,
		ID:       userObj.UserObjectId,
		Username: userObj.UserName,
	}

	token, err := r.service.GenerateJWT(userIdentity)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		helper.EncodeErrorResponse(w, err, "User created but error logging in. Try to refresh and try again", "500")
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

func (r resource) loginHandler(w http.ResponseWriter, req *http.Request) {
	var input dto.LoginRequestDto

	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		helper.EncodeErrorResponse(w, err, "Invalid Username or password", "400")
		return
	}

	if err := input.Validate(); err != nil {
		helper.EncodeErrorResponse(w, err, "Login failed, bad data!", "400")
		return
	}
	token, err := r.service.Login(req.Context(), input.Email, input.Password)
	if err != nil {
		helper.EncodeErrorResponse(w, err, "Invalid Username or password", "400")
		return
	}

	helper.EncodeSuccessResponse(w, token, "Login successful")
	return
}

func (r resource) getByIdHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	idk, _ := primitive.ObjectIDFromHex(id)

	category, _ := r.service.Get(req.Context(), idk)
	json.NewEncoder(w).Encode(category)
}
func (r resource) getUserPublicMetadata(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["userId"]
	idk, _ := primitive.ObjectIDFromHex(id)

	user, _ := r.service.GetUserProfileData(req.Context(), idk)

	helper.EncodeSuccessResponse(w, user, "User data retrieved.")
}

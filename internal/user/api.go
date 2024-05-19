package user

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/mux"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger) {
	res := resource{service, logger}
	r.HandleFunc("/api/v1/users/register", res.registerUserHandler).Methods("POST")

}

type CreateUserRequest struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Roles     []string `json:"roles"`
}

// Validate validates the CreateAlbumRequest fields.
func (m CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.FirstName, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.LastName, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Email, validation.Required, is.Email, validation.Length(6, 200)),
		validation.Field(&m.Password, validation.Required, validation.Length(6, 100)),

		//validation.Field(&a.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
	)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) registerUserHandler(w http.ResponseWriter, req *http.Request) {
	var input CreateUserRequest
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
	_, err = r.service.Create(req.Context(), input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

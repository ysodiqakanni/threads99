package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateNewUserRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"` // optional

	//FirstName string   `json:"first_name"`
	//LastName  string   `json:"last_name"`
	//Phone     string   `json:"phone"`    // Optional
	//Roles     []string `json:"roles"`
}

type CreateNewUserResponseDto struct {
	Token        string
	UserObjectId primitive.ObjectID
	UserId       string
	UserName     string
}

// Validate validates the CreateAlbumRequest fields.
func (m CreateNewUserRequestDto) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Email, validation.Required, is.Email, validation.Length(6, 200)),
		validation.Field(&m.Password, validation.Required, validation.Length(6, 100)),

		//validation.Field(&a.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
	)
}

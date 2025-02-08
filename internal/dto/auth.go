package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type LoginRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m LoginRequestDto) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Email, validation.Required, validation.Length(2, 128)),
		validation.Field(&m.Password, validation.Required, validation.Length(5, 128)),
	)
}

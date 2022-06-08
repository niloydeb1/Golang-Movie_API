package v1

import (
	"crypto/rsa"
	"errors"
	"github.com/niloydeb1/Golang-Movie_API/enums"
	"net/mail"
	"time"
)

// RsaKeys contains RSA keys.
type RsaKeys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// UserRegistrationDto dto that holds user registration info.
type UserRegistrationDto struct {
	ID                 string                 `json:"id" bson:"id"`
	FirstName          string                 `json:"first_name" bson:"first_name" `
	LastName           string                 `json:"last_name" bson:"last_name"`
	Email              string                 `json:"email" bson:"email" `
	Phone              string                 `json:"phone" bson:"phone"`
	Password           string                 `json:"password" bson:"password" `
	Status             enums.STATUS           `json:"status" bson:"status"`
	CreatedDate        time.Time              `json:"created_date" bson:"created_date"`
	UpdatedDate        time.Time              `json:"updated_date" bson:"updated_date"`
}

// Validate validates UserRegistrationDto data
func (u UserRegistrationDto) Validate() error {
	if u.ID == "" {
		return errors.New("user id is required")
	}
	if u.FirstName == "" {
		return errors.New("first name is required")
	}
	if u.LastName == "" {
		return errors.New("last name is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return err
	}
	return nil
}

// GetUserFromUserRegistrationDto converts User from UserRegistrationDto
func GetUserFromUserRegistrationDto(u UserRegistrationDto) User {
	user := User{
		ID:                 u.ID,
		FirstName:          u.FirstName,
		LastName:           u.LastName,
		Email:              u.Email,
		Phone:              u.Phone,
		Password:           u.Password,
		Status:             u.Status,
		CreatedDate:        u.CreatedDate,
		UpdatedDate:        u.UpdatedDate,
	}
	return user
}
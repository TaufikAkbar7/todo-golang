package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type VerifyUserRequest struct {
	Token string `validate:"required,max=100"`
}

type UserReponseLogin struct {
	Token string `json:"token"`
}

type UserCreateRequest struct {
	Username string `json:"username" form:"username" validate:"required,max=100"`
	Password string `json:"password" form:"password" validate:"required,max=100"`
}

type UserChangePasswordRequest struct {
	Password string `json:"password" form:"password"`
}

type UserChangePassword struct {
	Password string    `json:"password"`
	UpdateAt time.Time `json:"updated_at"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type UserResultFindUsername struct {
	Username string    `json:"username"`
	ID       uuid.UUID `json:"id"`
	Password string    `json:"-"`
}

type UserCustomClaims struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

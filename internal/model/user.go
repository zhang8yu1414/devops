package model

type UserCreateInput struct {
	Passport string `v:"required"`
	Password string `v:"required"`
	Nickname string
}

type UserSignInInput struct {
	Passport string
	Password string
}

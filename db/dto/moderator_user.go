package dto

import "pharmafinder/types"

type SessionTokenResponseDTO struct {
	Token    string `db:"-" json:"token"`
	ValidFor int64  `db:"-" json:"validFor"`
}

type AuthenticatedModeratorUserResponseDTO struct {
	ID                    types.UUID              `db:"id" json:"id"`
	Username              string                  `db:"username" json:"username"`
	Email                 string                  `db:"email" json:"email"`
	FirstName             string                  `db:"first_name" json:"firstName"`
	LastName              string                  `db:"last_name" json:"lastName"`
	RegistrationTimestamp types.Time              `db:"registration_timestamp" json:"registrationTs"`
	LastLoginTimestamp    types.Time              `db:"last_login_timestamp" json:"lastLoginTs"`
	Administrator         bool                    `db:"administrator" json:"administrator"`
	Session               SessionTokenResponseDTO `db:"-" json:"session"`
}

type ModeratorUserListResponseDTO struct {
	ID            types.UUID `db:"id" json:"id"`
	Username      string     `db:"username" json:"username"`
	Email         string     `db:"email" json:"email"`
	FirstName     string     `db:"first_name" json:"firstName"`
	LastName      string     `db:"last_name" json:"lastName"`
	Administrator bool       `db:"administrator" json:"administrator"`
}

type ModeratorUserRegistrationDTO struct {
	Username  string `json:"username" validate:"required,lte=32"`
	Email     string `json:"email" validate:"required,email,lte=64"`
	Password  string `json:"password" validate:"required,lte=72"`
	FirstName string `json:"firstName" validate:"required,lte=64"`
	LastName  string `json:"lastName" validate:"required,lte=64"`
}

type ModeratorUserLoginDTO struct {
	Username string `json:"username" validate:"required,lte=32"`
	Password string `json:"password" validate:"required,lte=72"`
}

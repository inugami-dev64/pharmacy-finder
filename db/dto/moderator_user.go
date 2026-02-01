package dto

import "pharmafinder/types"

type SessionTokenResponseDTO struct {
	Token    string `db:"-" json:"token"`
	ValidFor int64  `db:"-" json:"validFor"`
}

type ModeratorUserProfileDTO struct {
	ID                    types.UUID `db:"id" json:"id"`
	Username              string     `db:"username" json:"username"`
	Email                 string     `db:"email" json:"email"`
	FirstName             string     `db:"first_name" json:"firstName"`
	LastName              string     `db:"last_name" json:"lastName"`
	RegistrationTimestamp types.Time `db:"registration_timestamp" json:"registrationTs"`
	LastLoginTimestamp    types.Time `db:"last_login_timestamp" json:"lastLoginTs"`
	Administrator         bool       `db:"administrator" json:"administrator"`
}

type AuthenticatedModeratorUserResponseDTO struct {
	ModeratorUserProfileDTO
	Session SessionTokenResponseDTO `db:"-" json:"session"`
}

type ModeratorUserRegistrationDTO struct {
	Username  string `json:"username" validate:"required,lte=32"`
	Email     string `json:"email" validate:"required,email,lte=64"`
	Password  string `json:"password" validate:"required,lte=72"`
	FirstName string `json:"firstName" validate:"required,lte=64"`
	LastName  string `json:"lastName" validate:"required,lte=64"`
}

type AdminUserUpdateDTO struct {
	Email         string `json:"email" validate:"required,email,lte=64"`
	FirstName     string `json:"firstName" validate:"required,lte=64"`
	LastName      string `json:"lastName" validate:"required,lte=64"`
	Password      string `json:"password" validate:"lte=72"`
	Administrator *bool  `json:"administrator"`
}

type ModeratorUserUpdateDTO struct {
	AdminUserUpdateDTO
	CurrentPassword string `json:"currentPassword" validate:"required,lte=72"`
}

type ModeratorUserDeletionDTO struct {
	CurrentPassword string `json:"currentPassword" validate:"required,lte=72"`
}

type ModeratorUserLoginDTO struct {
	Username string `json:"username" validate:"required,lte=32"`
	Password string `json:"password" validate:"required,lte=72"`
}

package dto

import "pharmafinder/types"

type AuthenticatedModeratorUserResponseDTO struct {
	ID                    types.UUID `db:"id" json:"id"`
	Username              string     `db:"username" json:"username"`
	Email                 string     `db:"email" json:"email"`
	FirstName             string     `db:"first_name" json:"firstName"`
	LastName              string     `db:"last_name" json:"lastName"`
	RegistrationTimestamp types.Time `db:"registration_timestamp" json:"registrationTs"`
	LastLoginTimestamp    types.Time `db:"last_login_timestamp" json:"lastLoginTs"`
	Administrator         bool       `db:"administrator" json:"administrator"`
	Session               struct {
		Token    string `db:"-" json:"token"`
		ValidFor int64  `db:"-" json:"validFor"`
	} `db:"-" json:"session"`
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
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type ModeratorUserLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

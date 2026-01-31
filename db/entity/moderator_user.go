package entity

import "pharmafinder/types"

type ModeratorUser struct {
	ID                    types.UUID `db:"id"`
	Username              string     `db:"username"`
	Email                 string     `db:"email"`
	Password              string     `db:"password"`
	FirstName             string     `db:"first_name"`
	LastName              string     `db:"last_name"`
	RegistrationTimestamp types.Time `db:"registration_timestamp"`
	LastLoginTimestamp    types.Time `db:"last_login_timestamp"`
	Administrator         bool       `db:"administrator"`
}

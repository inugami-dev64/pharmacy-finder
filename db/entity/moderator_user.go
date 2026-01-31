package entity

import "pharmafinder/types"

type ModeratorUser struct {
	ID            types.UUID `db:"id"`
	Username      string     `db:"username"`
	Password      string     `db:"password"`
	FirstName     string     `db:"first_name"`
	LastName      string     `db:"last_name"`
	Administrator bool       `db:"administrator"`
}

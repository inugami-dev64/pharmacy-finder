package db

import (
	"pharmafinder/db/dto"
	"pharmafinder/db/entity"
	"pharmafinder/types"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ModeratorUserRepository interface {
	FindUserByID(id types.UUID) Query[entity.ModeratorUser]
	FindUserByUsernameOrEmail(usernameOrEmail string) Query[entity.ModeratorUser]
	FindAllUsers() Query[dto.AuthenticatedModeratorUserResponseDTO]
	HasAdministrator() Query[bool]

	Store(user *entity.ModeratorUser) error
	Trx(conn any) ModeratorUserRepository
}

type ModeratorUserRepositorySQLX struct {
	conn *sqlx.DB
}

func ProvideModeratorUserRepository(conn *sqlx.DB) ModeratorUserRepository {
	return ModeratorUserRepositorySQLX{conn: conn}
}

func (repo ModeratorUserRepositorySQLX) FindUserByID(id types.UUID) Query[entity.ModeratorUser] {
	q := `
	SELECT
		*
	FROM
		moderator_users mu
	WHERE
		mu.id = $1
	`

	args := []interface{}{id}

	return &SQLXQuery[entity.ModeratorUser]{
		uniqueKey: "id",
		key:       "username",
		trx:       repo.conn,
		q:         q,
		args:      args,
	}
}

func (repo ModeratorUserRepositorySQLX) FindUserByUsernameOrEmail(usernameOrEmail string) Query[entity.ModeratorUser] {
	q := `
	SELECT
		*
	FROM
		moderator_users mu
	WHERE
		mu.username = $1
	OR
		mu.email = $2
	`

	args := []interface{}{usernameOrEmail, usernameOrEmail}

	return &SQLXQuery[entity.ModeratorUser]{
		uniqueKey: "id",
		key:       "username",
		trx:       repo.conn,
		q:         q,
		args:      args,
	}
}

func (repo ModeratorUserRepositorySQLX) FindAllUsers() Query[dto.AuthenticatedModeratorUserResponseDTO] {
	q := `
	SELECT
		mu.id,
		mu.username,
		mu.email,
		mu.first_name,
		mu.last_name,
		mu.registration_timestamp,
		mu.last_login_timestamp,
		mu.administrator
	FROM
		moderator_users mu
	`

	return &SQLXQuery[dto.AuthenticatedModeratorUserResponseDTO]{
		uniqueKey: "id",
		key:       "username",
		trx:       repo.conn,
		q:         q,
	}
}

func (repo ModeratorUserRepositorySQLX) HasAdministrator() Query[bool] {
	q := `
	SELECT
		COUNT(*) > 0
	FROM moderator_users mu
	WHERE
		mu.administrator
	`

	return &SQLXQuery[bool]{
		trx: repo.conn,
		q:   q,
	}
}

func (repo ModeratorUserRepositorySQLX) Store(user *entity.ModeratorUser) error {
	var rows *sqlx.Rows
	var err error
	if user.ID != types.UUID(uuid.Nil) {
		rows, err = repo.conn.NamedQuery(
			`UPDATE moderator_users SET
				username = :username,
				email = :email,
				"password" = :password,
				first_name = :first_name,
				last_name = :last_name,
				registration_timestamp = :registration_timestamp,
				last_login_timestamp = :last_login_timestamp,
				administrator = :administrator
			WHERE
				id = :id
			RETURNING *
			`, user)
	} else {
		rows, err = repo.conn.NamedQuery(
			`INSERT INTO moderator_users (username,email,"password",first_name,last_name,registration_timestamp,last_login_timestamp,administrator)
				VALUES (:username,:email,:password,:first_name,:last_name,:registration_timestamp,:last_login_timestamp,:administrator)
			RETURNING *`,
			user)
	}

	if err != nil {
		return err
	}

	for rows.Next() {
		_ = rows.StructScan(user)
	}

	return nil
}

func (repo ModeratorUserRepositorySQLX) Trx(conn any) ModeratorUserRepository {
	return ModeratorUserRepositorySQLX{conn: conn.(*sqlx.DB)}
}

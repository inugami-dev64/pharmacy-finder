package db

import (
	"pharmafinder/db/entity"

	"github.com/jmoiron/sqlx"
)

type PharmacyReviewRepository interface {
	FindReviewForPharmacy(id int64) Query[entity.PharmacyReview]
	FindReviewByID(pharmaID int64, reviewID int64) Query[entity.PharmacyReview]
	Store(review *entity.PharmacyReview) error
	Trx(conn any) PharmacyReviewRepository
}

type PharmacyReviewRepositorySQLX struct {
	conn *sqlx.DB
}

func ProvidePharmacyReviewRepository(conn *sqlx.DB) PharmacyReviewRepository {
	return PharmacyReviewRepositorySQLX{conn: conn}
}

func (repo PharmacyReviewRepositorySQLX) FindReviewForPharmacy(id int64) Query[entity.PharmacyReview] {
	q := `
	SELECT
		*
	FROM
		pharmacy_reviews pr
	WHERE
		pr.pharmacy_id = $1
	`

	args := []interface{}{id}

	return &SQLXQuery[entity.PharmacyReview]{
		uniqueKey: "id",
		key:       "updated_at",
		trx:       repo.conn,
		q:         q,
		args:      args,
	}
}

func (repo PharmacyReviewRepositorySQLX) FindReviewByID(pharmaID int64, reviewID int64) Query[entity.PharmacyReview] {
	q := `
	SELECT
		*
	FROM
		pharmacy_reviews pr
	WHERE
		pr.pharmacy_id = $1
	AND
		pr.id = $2
	`

	args := []interface{}{pharmaID, reviewID}

	return &SQLXQuery[entity.PharmacyReview]{
		uniqueKey: "id",
		key:       "updated_at",
		trx:       repo.conn,
		q:         q,
		args:      args,
	}
}

func (repo PharmacyReviewRepositorySQLX) Store(review *entity.PharmacyReview) error {
	if review.ID != 0 {
		_, err := repo.conn.NamedExec(
			`UPDATE pharmacy_reviews SET
				pharmacy_id = :pharmacy_id,
				prescription_type = :prescription_type,
				stars = :stars,
				hrt_kind = :hrt_kind,
				nationality = :nationality,
				review = :review,
				created_at = :created_at,
				updated_at = :updated_at
			WHERE
				id = :id
			`, review)
		return err
	}

	rows, err := repo.conn.NamedQuery(
		`INSERT INTO pharmacy_reviews (pharmacy_id,prescription_type,stars,hrt_kind,nationality,review,created_at,updated_at,modification_code)
			VALUES (:pharmacy_id,:prescription_type,:stars,:hrt_kind,:nationality,:review,:created_at,:updated_at,:modification_code)
		RETURNING *`,
		review)

	if err != nil {
		return err
	}

	for rows.Next() {
		_ = rows.StructScan(review)
	}

	return nil
}

func (repo PharmacyReviewRepositorySQLX) Trx(conn any) PharmacyReviewRepository {
	return PharmacyReviewRepositorySQLX{conn: conn.(*sqlx.DB)}
}

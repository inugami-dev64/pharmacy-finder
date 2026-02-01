package db

import (
	"pharmafinder/db/dto"
	"pharmafinder/db/entity"
	"strings"

	"github.com/jmoiron/sqlx"
)

type PharmacyReviewRepository interface {
	FindReviewForPharmacy(id int64) Query[entity.PharmacyReview]
	FindReviewByID(pharmaID int64, reviewID int64) Query[entity.PharmacyReview]
	FindReviewsForModeration(showUnmoderated bool, showModerated bool) Query[dto.ModerationPharmacyReview]
	Store(review *entity.PharmacyReview) error
	Delete(id int64) Query[entity.PharmacyReview]
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

func (repo PharmacyReviewRepositorySQLX) FindReviewsForModeration(showUnmoderated bool, showModerated bool) Query[dto.ModerationPharmacyReview] {
	q := `
	SELECT
		pr.id,
		pr.prescription_type,
		pr.stars,
		pr.hrt_kind,
		pr.nationality,
		pr.review,
		pr.created_at,
		pr.updated_at,
		pr.pharmacy_id,
		CASE
			WHEN ROUND(mr.avg_result) = 1 THEN 'APPROVED'::comment_review_result_t
			WHEN ROUND(mr.avg_result) = 2 THEN 'OTHER'::comment_review_result_t
			WHEN ROUND(mr.avg_result) = 3 THEN 'OFFENSIVE'::comment_review_result_t
			WHEN ROUND(mr.avg_result) = 4 THEN 'PERSONAL_ATTACK'::comment_review_result_t
			ELSE 'NONE'::comment_review_result_t
		END AS "result",
		COALESCE(mr.avg_marked_for_deletion, 0) >= 0.5 AS marked_for_deletion,
		mr.reviewed_at
	FROM
		pharmacy_reviews pr
	LEFT JOIN (
		SELECT
			cr.comment_id,
			AVG(
				CASE
					WHEN cr."result" = 'APPROVED' THEN 1
					WHEN cr."result" = 'OTHER' THEN 2
					WHEN cr."result" = 'OFFENSIVE' THEN 3
					WHEN cr."result" = 'PERSONAL_ATTACK' THEN 4
				END
			) AS avg_result,
			AVG(CAST(cr.marked_for_deletion AS INT)) AS avg_marked_for_deletion,
			MAX(cr."reviewed_at") AS reviewed_at
		FROM
			comment_reviews cr
		GROUP BY
			cr.comment_id
	) mr
	ON
		mr.comment_id = pr.id
	`

	whereClauses := []string{}
	if showUnmoderated {
		whereClauses = append(whereClauses, "mr.comment_id IS NULL")
	}
	if showModerated {
		whereClauses = append(whereClauses, "mr.comment_id IS NOT NULL")
	}

	whereClause := strings.Join(whereClauses, " OR ")
	if whereClause != "" {
		q += " WHERE " + whereClause
	}

	return &SQLXQuery[dto.ModerationPharmacyReview]{
		uniqueKey: "id",
		key:       "updated_at",
		trx:       repo.conn,
		q:         q,
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

// Warning: potentially destructive action
func (repo PharmacyReviewRepositorySQLX) Delete(id int64) Query[entity.PharmacyReview] {
	q := `
	DELETE FROM pharmacy_reviews
	WHERE id = $1
	RETURNING *
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

func (repo PharmacyReviewRepositorySQLX) Trx(conn any) PharmacyReviewRepository {
	return PharmacyReviewRepositorySQLX{conn: conn.(*sqlx.DB)}
}

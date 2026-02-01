package db

import (
	"pharmafinder/db/dto"
	"pharmafinder/db/entity"

	"github.com/jmoiron/sqlx"
)

type CommentModerationRepository interface {
	FindCommentModerationsForReview(reviewID int64) Query[dto.CommentReviewResultDTO]
	FindCommentModerationByReviewAndID(reviewID int64, modID int64) Query[entity.CommentReview]

	Store(review *entity.CommentReview) error
	Delete(modID int64) Query[entity.CommentReview]
	Trx(conn any) CommentModerationRepository
}

type CommentModerationRepositorySQLX struct {
	conn *sqlx.DB
}

func ProvideModerationRepository(conn *sqlx.DB) CommentModerationRepository {
	return CommentModerationRepositorySQLX{conn: conn}
}

func (repo CommentModerationRepositorySQLX) FindCommentModerationsForReview(reviewID int64) Query[dto.CommentReviewResultDTO] {
	q := `
	SELECT
		cr.*,
		mu.username AS moderator_username
	FROM
		comment_reviews cr
	INNER JOIN
		moderator_users mu
	ON
		cr.moderator_id = mu.id
	WHERE
		cr.comment_id = $1
	ORDER BY
		reviewed_at DESC`

	args := []interface{}{reviewID}

	return &SQLXQuery[dto.CommentReviewResultDTO]{
		uniqueKey: "id",
		key:       "reviewed_at",
		q:         q,
		args:      args,
	}
}

func (repo CommentModerationRepositorySQLX) FindCommentModerationByReviewAndID(reviewID int64, modID int64) Query[entity.CommentReview] {
	q := `
	SELECT
		*
	FROM
		comment_reviews cr
	WHERE
		cr.comment_id = $1
	AND
		cr.id = $2`

	args := []interface{}{reviewID, modID}

	return &SQLXQuery[entity.CommentReview]{
		uniqueKey: "id",
		key:       "reviewed_at",
		q:         q,
		args:      args,
	}
}

func (repo CommentModerationRepositorySQLX) Store(review *entity.CommentReview) error {
	var rows *sqlx.Rows
	var err error
	if review.ID != 0 {
		rows, err = repo.conn.NamedQuery(
			`UPDATE comment_reviews SET
				result = :result,
				mod_comment = :mod_comment,
				marked_for_deletion = :marked_for_deletion,
				reviewed_at = :reviewed_at
			WHERE
				id = :id
			RETURNING *`, review)
	} else {
		rows, err = repo.conn.NamedQuery(
			`INSERT INTO comment_reviews(comment_id,moderator_id,result,mod_comment,marked_for_deletion,reviewed_at)
				VALUES(:comment_id,:moderator_id,:result,:mod_comment,:marked_for_deletion,:reviewed_at)
			RETURNING *`, review)
	}

	if err != nil {
		return err
	}

	for rows.Next() {
		_ = rows.StructScan(review)
	}

	return nil
}

func (repo CommentModerationRepositorySQLX) Delete(modID int64) Query[entity.CommentReview] {
	q := `
	DELETE
	FROM
		comment_reviews
	WHERE
		id = $1
	RETURNING *`

	args := []interface{}{modID}

	return &SQLXQuery[entity.CommentReview]{
		uniqueKey: "id",
		key:       "reviewed_at",
		q:         q,
		args:      args,
	}
}

func (repo CommentModerationRepositorySQLX) Trx(conn any) CommentModerationRepository {
	return CommentModerationRepositorySQLX{conn: conn.(*sqlx.DB)}
}

package entity

import "pharmafinder/types"

type CommentReviewResult string

const (
	COM_REVIEW_APPROVED        CommentReviewResult = CommentReviewResult("APPROVED")
	COM_REVIEW_PERSONAL_ATTACK                     = CommentReviewResult("PERSONAL_ATTACK")
	COM_REVIEW_OFFENSIVE                           = CommentReviewResult("OFFENSIVE")
	COM_REVIEW_OTHER                               = CommentReviewResult("OTHER")
)

type CommentReview struct {
	ID                int64               `db:"id"`
	CommentID         int64               `db:"comment_id"`
	ModeratorID       types.UUID          `db:"moderator_id"`
	Result            CommentReviewResult `db:"result"`
	ModeratorComment  *string             `db:"mod_comment"`
	MarkedForDeletion bool                `db:"marked_for_deletion"`
	ReviewedAt        types.Time          `db:"reviewed_at"`
}

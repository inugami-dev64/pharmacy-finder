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
	ID                int64               `db:"id" json:"id"`
	CommentID         int64               `db:"comment_id" json:"commentId"`
	ModeratorID       types.UUID          `db:"moderator_id" json:"moderatorId"`
	Result            CommentReviewResult `db:"result" json:"result"`
	ModeratorComment  *string             `db:"mod_comment" json:"moderatorComment"`
	MarkedForDeletion bool                `db:"marked_for_deletion" json:"markedForDeletion"`
	ReviewedAt        types.Time          `db:"reviewed_at" json:"reviewedAt"`
}

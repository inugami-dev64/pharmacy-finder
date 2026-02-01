package dto

import "pharmafinder/db/entity"

type CommentReviewResultDTO struct {
	entity.CommentReview
	ModeratorName string `db:"moderator_username" json:"moderatorUsername"`
}

type CommentReviewModificationDTO struct {
	Result            entity.CommentReviewResult `json:"result" validate:"required,oneof=APPROVED PERSONAL_ATTACK OFFENSIVE OTHER"`
	ModeratorComment  *string                    `json:"comment"`
	MarkedForDeletion bool                       `json:"markedForDeletion"`
}

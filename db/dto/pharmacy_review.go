package dto

import "pharmafinder/types"

type PharmacyReviewCreationDTO struct {
	PrescriptionType  string  `json:"prescriptionType" validate:"required,oneof=Imago GenderGP National"`
	Stars             int     `json:"stars" validate:"required"`
	HRTKind           string  `json:"hrtKind" validate:"required,oneof=t e"`
	Nationality       *string `json:"nationality" validate:"iso3166_1_alpha2"`
	Review            *string `json:"review" validate:"lte=1024"`
	RecaptchaResponse string  `json:"__gRecaptchaResponse"`
}

type PharmacyReviewDeletionDTO struct {
	ModificationCode  string `json:"modCode" validate:"required,lte=16"`
	RecaptchaResponse string `json:"__gRecaptchaResponse"`
}

type PharmacyReviewsetResultDTO struct {
	ID               int64      `db:"id" json:"id"`
	PrescriptionType string     `db:"prescription_type" json:"prescriptionType"`
	Stars            int        `db:"stars" json:"stars"`
	HRTKind          string     `db:"hrt_kind" json:"hrtKind"`
	Nationality      *string    `db:"nationality" json:"nationality"`
	Review           *string    `db:"review" json:"review"`
	CreatedAt        types.Time `db:"created_at" json:"createdAt"`
	UpdatedAt        types.Time `db:"updated_at" json:"updatedAt"`
}

type ModerationPharmacyReview struct {
	PharmacyReviewsetResultDTO
	PharmacyID          int64       `db:"pharmacy_id" json:"pharmacyId"`
	CommentReviewResult string      `db:"result" json:"commentReviewResult"`
	MarkedForDeletion   bool        `db:"marked_for_deletion" json:"markedForDeletion"`
	ReviewedAt          *types.Time `db:"reviewed_at" json:"reviewedAt,omitempty"`
}

type PharmacyReviewModificationDTO struct {
	PharmacyReviewCreationDTO
	ModificationCode string `json:"modCode" validate:"required,lte=16"`
}

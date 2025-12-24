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
	ID               int64      `json:"id"`
	PrescriptionType string     `json:"prescriptionType"`
	Stars            int        `json:"stars"`
	HRTKind          string     `json:"hrtKind"`
	Nationality      *string    `json:"nationality"`
	Review           *string    `json:"review"`
	CreatedAt        types.Time `json:"createdAt"`
	UpdatedAt        types.Time `json:"updatedAt"`
}

type PharmacyReviewModificationDTO struct {
	PharmacyReviewCreationDTO
	ModificationCode string `json:"modCode" validate:"required,lte=16"`
}

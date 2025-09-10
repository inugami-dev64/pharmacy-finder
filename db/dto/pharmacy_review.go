package dto

type PharmacyReviewCreationDTO struct {
	PrescriptionType string  `json:"prescriptionType" validate:"required"`
	Stars            int     `json:"stars" validate:"required"`
	HRTKind          string  `json:"hrtKind" validate:"required"`
	Nationality      *string `json:"nationality" validate:"iso3166_1_alpha2"`
	Review           *string `json:"review" validate:"lte=1024"`
}

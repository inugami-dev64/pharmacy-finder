package dto

type PharmacyRatingDTO struct {
	ID      int64   `db:"id" json:"id" `
	Stars   float32 `db:"stars" json:"stars"`
	HRTKind *string `db:"hrt_kind" json:"hrtKind"`
}

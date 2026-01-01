package dto

type PharmacyRatingDTO struct {
	ID      int64   `db:"id" json:"id" `
	Stars   float32 `db:"stars" json:"stars"`
	HRTKind *string `db:"hrt_kind" json:"hrtKind"`
}

type PharmacyTierRatingDTO struct {
	ID         int64   `db:"id" json:"id"`
	Name       string  `db:"name" json:"name"`
	AvgRating  float64 `db:"avg_rating" json:"avgRating"`
	AvgERating float64 `db:"avg_e_rating" json:"avgERating"`
	AvgTRating float64 `db:"avg_t_rating" json:"avgTRating"`
}

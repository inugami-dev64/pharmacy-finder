package dto

type PharmacyRating struct {
	OverallRating      *float32 `json:"rating"`
	AcceptanceRating   *float32 `json:"acceptance"`
	EstrogenRating     *float32 `json:"e"`
	TestosteroneRating *float32 `json:"t"`
}

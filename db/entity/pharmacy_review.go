package entity

import "pharmafinder/types"

type PrescriptionType string

const (
	PRESCRIPTION_IMAGO    PrescriptionType = PrescriptionType("Imago")
	PRESCRIPTION_GENDERGP                  = PrescriptionType("GenderGP")
	PRESCRIPTION_NATIONAL                  = PrescriptionType("National")
)

type HRTKind string

const (
	HRT_KIND_ESTROGEN_BASED     HRTKind = "e"
	HRT_KIND_TESTOSTERONE_BASED HRTKind = "t"
)

type PharmacyReview struct {
	ID               int64      `db:"id" json:"id"`
	PharmacyID       int64      `db:"pharmacy_id" json:"pharmacyId"`
	PrescriptionType string     `db:"type" json:"prescriptionType"`
	Stars            int        `db:"stars" json:"stars"`
	HRTKind          string     `db:"hrt_kind" json:"hrtKind"`
	Nationality      *string    `db:"nationality" json:"nationality"`
	Review           *string    `db:"review" json:"review"`
	CreatedAt        types.Time `db:"created_at" json:"createdAt"`
	UpdatedAt        types.Time `db:"updated_at" json:"updatedAt"`
	ModificationCode string     `db:"mod_code" json:"modCode"`
}

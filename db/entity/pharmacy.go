package entity

import "pharmafinder/types"

type PharmacyChain string

const (
	CHAIN_APOTHEKA        PharmacyChain = PharmacyChain("Apotheka")
	CHAIN_SUDAMEAPTEEK                  = PharmacyChain("Südameapteek")
	CHAIN_BENU                          = PharmacyChain("Benu")
	CHAIN_EUROAPTEEK                    = PharmacyChain("Euroapteek")
	CHAIN_KALAMAJA_APTEEK               = PharmacyChain("Kalamaja")
)

type Pharmacy struct {
	ID          int64      `db:"id" json:"id"`
	PharmacyID  int64      `db:"pharmacy_id" json:"-"`
	Chain       string     `db:"chain" json:"chain"`
	Name        string     `db:"name" json:"name"`
	Address     string     `db:"address" json:"address"`
	City        string     `db:"city" json:"city"`
	County      string     `db:"county" json:"county"`
	PostalCode  string     `db:"postal_code" json:"postalCode"`
	Email       string     `db:"email" json:"email"`
	PhoneNumber string     `db:"phone_number" json:"phoneNumber"`
	ModTime     types.Time `db:"mod_time" json:"-"`
	Latitude    float32    `db:"latitude" json:"lat"`
	Longitude   float32    `db:"longitude" json:"lng"`
}

package entity

type Pharmacy struct {
	ID          int64   `db:"id" json:"id"`
	Chain       string  `db:"chain" json:"chain"`
	Name        string  `db:"name" json:"name"`
	Address     string  `db:"address" json:"address"`
	City        string  `db:"city" json:"city"`
	County      string  `db:"county" json:"county"`
	PostalCode  string  `db:"postal_code" json:"postalCode"`
	PhoneNumber string  `db:"phone_number" json:"phoneNumber"`
	Latitude    float32 `db:"latitude" json:"lat"`
	Longitude   float32 `db:"longitude" json:"lng"`
}

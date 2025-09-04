package db

import (
	"pharmafinder/db/entity"
	"pharmafinder/types"

	"github.com/jmoiron/sqlx"
)

type PharmacyRepository interface {
	FindPharmaciesInCoordinateBounds(sw types.Point, ne types.Point) Query[entity.Pharmacy]
	FindPharmaciesByChain(chain entity.PharmacyChain) Query[entity.Pharmacy]
	StoreAll(pharmacies []entity.Pharmacy) error
	Trx(conn any) PharmacyRepository
}

type PharmacyRepositorySQLX struct {
	conn *sqlx.DB
}

func ProvidePharmacyRepository(conn *sqlx.DB) PharmacyRepository {
	return PharmacyRepositorySQLX{conn: conn}
}

func (repo PharmacyRepositorySQLX) FindPharmaciesInCoordinateBounds(sw types.Point, ne types.Point) Query[entity.Pharmacy] {
	q := `
	SELECT
		*
	FROM
		pharmacies p
	WHERE
		p.latitude >= $1
	AND
		p.longitude >= $2
	AND
		p.latitude <= $3
	AND
		p.longitude <= $4
	`

	args := []interface{}{sw.Lat, sw.Lng, ne.Lat, ne.Lng}

	return &SQLXQuery[entity.Pharmacy]{
		uniqueKey: "id",
		key:       "id",
		trx:       repo.conn,
		q:         q,
		args:      args,
	}
}

func (repo PharmacyRepositorySQLX) FindPharmaciesByChain(chain entity.PharmacyChain) Query[entity.Pharmacy] {
	q := `
	SELECT
		*
	FROM
		pharmacies p
	WHERE
		chain = ?
	`

	args := []interface{}{string(chain)}
	return &SQLXQuery[entity.Pharmacy]{
		uniqueKey: "id",
		key:       "id",
		trx:       repo.conn,
		q:         q,
		args:      args,
	}
}

func (repo PharmacyRepositorySQLX) StoreAll(pharmacies []entity.Pharmacy) error {
	// Separate entities which shall be inserted
	// and entities which shall be updated
	_, err := repo.conn.NamedExec(
		`INSERT INTO pharmacies (chain,"name","address",city,county,postal_code,phone_number,latitude,longitude)
			VALUES (:chain,:name,:address,:city,:county,:postal_code,:phone_number,:latitude,:longitude)
		ON CONFLICT DO UPDATE
		`, pharmacies)

	return err
}

func (repo PharmacyRepositorySQLX) Trx(conn any) PharmacyRepository {
	return PharmacyRepositorySQLX{conn: conn.(*sqlx.DB)}
}

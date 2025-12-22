package db

import (
	"pharmafinder/db/dto"
	"pharmafinder/db/entity"
	"pharmafinder/types"

	"github.com/jmoiron/sqlx"
)

type PharmacyRepository interface {
	FindPharmaciesInCoordinateBounds(sw types.Point, ne types.Point) Query[entity.Pharmacy]
	FindPharmaciesByChain(chain entity.PharmacyChain) Query[entity.Pharmacy]
	FindPharmacyRatingsByID(id int64) Query[dto.PharmacyRatingDTO]
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
		p.chain = $1
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

func (repo PharmacyRepositorySQLX) FindPharmacyRatingsByID(id int64) Query[dto.PharmacyRatingDTO] {
	q := `SELECT * FROM find_pharmacy_ratings($1)`

	args := []interface{}{id}
	return &SQLXQuery[dto.PharmacyRatingDTO]{
		uniqueKey: "hrt_kind",
		key:       "id",
		trx:       repo.conn,
		q:         q,
		args:      args,
	}
}

func (repo PharmacyRepositorySQLX) StoreAll(pharmacies []entity.Pharmacy) error {
	// Separate entities which shall be inserted
	// and entities which shall be updated
	toInsert := make([]entity.Pharmacy, 0)
	for _, entity := range pharmacies {
		if entity.ID != 0 {
			_, err := repo.conn.NamedExec(
				`UPDATE pharmacies SET
					pharmacy_id = :pharmacy_id,
					chain = :chain,
					name = :name,
					"address" = :address,
					city = :city,
					county = :county,
					postal_code = :postal_code,
					email = :email,
					phone_number = :phone_number,
					mod_time = :mod_time,
					latitude = :latitude,
					longitude = :longitude
				WHERE
					id = :id
				`, entity)
			if err != nil {
				return err
			}
		} else {
			toInsert = append(toInsert, entity)
		}
	}

	if len(toInsert) > 0 {
		_, err := repo.conn.NamedExec(
			`INSERT INTO pharmacies (pharmacy_id,chain,"name","address",city,county,postal_code,email,phone_number,mod_time,latitude,longitude)
				VALUES (:pharmacy_id,:chain,:name,:address,:city,:county,:postal_code,:email,:phone_number,:mod_time,:latitude,:longitude)`,
			toInsert)
		return err
	}

	return nil
}

func (repo PharmacyRepositorySQLX) Trx(conn any) PharmacyRepository {
	return PharmacyRepositorySQLX{conn: conn.(*sqlx.DB)}
}

package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"shinagawa-shelter/pkg/config"
	"shinagawa-shelter/pkg/model"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	QUERY_INSERT = `INSERT INTO shelters (
			facility_name,
			facility_name_en,
			category,
			prefecture,
			city,
			administrative_district,
			address,
			address_kana,
			telephone_number,
			target_district,
			target_district_description,
			latitude,
			longitude,
			above_sea_level,
			facility_capacity,
			facility_height,
			flood,
			landslide,
			storm_surge,
			earthquake,
			tsunami,
			inundation,
			fire,
			eruption,
			description,
			description_en,
			standard_area_code,
			created_at
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13,
			$14,
			$15,
			$16,
			$17,
			$18,
			$19,
			$20,
			$21,
			$22,
			$23,
			$24,
			$25,
			$26,
			$27,
			$28
		) RETURNING id`

	QUERY_UPDATE = `UPDATE shelters SET
			facility_name_en = $1,
			category = $2,
			prefecture = $3,
			city = $4,
			administrative_district = $5,
			address = $6,
			address_kana = $7,
			telephone_number = $8,
			target_district = $9,
			target_district_description = $10,
			latitude = $11,
			longitude = $12,
			above_sea_level = $13,
			facility_capacity = $14,
			facility_height = $15,
			flood = $16,
			landslide = $17,
			storm_surge = $18,
			earthquake = $19,
			tsunami = $20,
			inundation = $21,
			fire = $22,
			eruption = $23,
			description = $24,
			description_en = $25,
			standard_area_code = $26,
			updated_at = $27
		WHERE
			facility_name = $28
	`

	QUERY_DELETE = `DELETE FROM shelters
		WHERE
			facility_name = $1
	`

	QUERY_SELECT = `SELECT
			id,
			facility_name,
			facility_name_en,
			category,
			prefecture,
			city,
			administrative_district,
			address,
			address_kana,
			telephone_number,
			target_district,
			target_district_description,
			latitude,
			longitude,
			above_sea_level,
			facility_capacity,
			facility_height,
			flood,
			landslide,
			storm_surge,
			earthquake,
			tsunami,
			inundation,
			fire,
			eruption,
			description,
			description_en,
			standard_area_code,
			created_at,
			updated_at
		FROM
			shelters
		WHERE
			facility_name = $1
	`

	QUERY_SELECT_ALL = `SELECT
			id,
			facility_name,
			facility_name_en,
			category,
			prefecture,
			city,
			administrative_district,
			address,
			address_kana,
			telephone_number,
			target_district,
			target_district_description,
			latitude,
			longitude,
			above_sea_level,
			facility_capacity,
			facility_height,
			flood,
			landslide,
			storm_surge,
			earthquake,
			tsunami,
			inundation,
			fire,
			eruption,
			description,
			description_en,
			standard_area_code,
			created_at,
			updated_at
		FROM
			shelters
		WHERE
			0 = 0
	`
)

var (
	db *sql.DB

	// defaultPage Default Page number
	defaultPage uint64 = 1

	// defaultPageSize Default Page size
	defaultPageSize uint64 = 10
)

type Condition struct {
	Category  *string
	Latitude  *float64
	Longitude *float64
	Page      uint64
	PageSize  uint64
}

func init() {
	var err error
	db, err = sql.Open("pgx", config.DataSourceName)
	if err != nil {
		log.Fatal("failed to open database", err)
	}
	log.Println("connected to database")
}

func NewCondition(c map[string]string) *Condition {
	if len(c) == 0 {
		return nil
	}

	category := c["category"]
	lat, err := strconv.ParseFloat(c["lat"], 64)
	if lat < 1 || err != nil {
		lat = 0
	}
	lon, err := strconv.ParseFloat(c["lon"], 64)
	if lon < 1 || err != nil {
		lon = 0
	}
	page, err := strconv.ParseUint(c["page"], 10, 64)
	if err != nil {
		page = defaultPage
	}
	size, err := strconv.ParseUint(c["size"], 10, 64)
	if err != nil {
		size = defaultPageSize
	}

	return &Condition{
		Category:  &category,
		Latitude:  &lat,
		Longitude: &lon,
		Page:      page,
		PageSize:  size,
	}
}

func (c *Condition) Build() string {
	var query string

	// category
	if len(*c.Category) > 0 {
		query += fmt.Sprintf(" AND category::jsonb ? '%s' ", *c.Category)
	}
	// order by distance
	if *c.Latitude > 0 && *c.Longitude > 0 {
		query += fmt.Sprintf(" ORDER BY POINT(latitude, longitude) <-> POINT(%f, %f) ", *c.Latitude, *c.Longitude)
	}
	// pagenation
	query += fmt.Sprintf(" LIMIT %d OFFSET %d ", c.PageSize, (c.Page-1)*c.PageSize)

	return query
}

func FindShelters(c *Condition) ([]*model.Shelter, error) {
	query := QUERY_SELECT_ALL

	if c != nil {
		// build where clause
		query += c.Build()
	} else {
		// default order by id
		query += " ORDER BY id"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	shelters := make([]*model.Shelter, 0)

	for rows.Next() {
		s := &model.Shelter{}
		var categoryJSON string

		err := rows.Scan(
			&s.ID,
			&s.FacilityName,
			&s.EnglishFacilityName,
			&categoryJSON,
			&s.Prefecture,
			&s.City,
			&s.AdministrativeDistrict,
			&s.Address,
			&s.AddressKana,
			&s.TelephoneNumber,
			&s.TargetDistrict,
			&s.TargetDistrictDescription,
			&s.Latitude,
			&s.Longitude,
			&s.AboveSeaLevel,
			&s.FacilityCapacity,
			&s.FacilityHeight,
			&s.Flood,
			&s.Landslide,
			&s.StormSurge,
			&s.Earthquake,
			&s.Tsunami,
			&s.Inundation,
			&s.Fire,
			&s.Eruption,
			&s.Description,
			&s.EnglishDescription,
			&s.StandardAreaCode,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(categoryJSON), &s.Category)
		if err != nil {
			return nil, err
		}

		shelters = append(shelters, s)
	}
	if err != nil {
		return nil, err
	}

	return shelters, nil
}

func FindShelter(facilityName string) (*model.Shelter, error) {
	s := &model.Shelter{}
	var categoryJSON string

	err := db.QueryRow(QUERY_SELECT, facilityName).Scan(
		&s.ID,
		&s.FacilityName,
		&s.EnglishFacilityName,
		&categoryJSON,
		&s.Prefecture,
		&s.City,
		&s.AdministrativeDistrict,
		&s.Address,
		&s.AddressKana,
		&s.TelephoneNumber,
		&s.TargetDistrict,
		&s.TargetDistrictDescription,
		&s.Latitude,
		&s.Longitude,
		&s.AboveSeaLevel,
		&s.FacilityCapacity,
		&s.FacilityHeight,
		&s.Flood,
		&s.Landslide,
		&s.StormSurge,
		&s.Earthquake,
		&s.Tsunami,
		&s.Inundation,
		&s.Fire,
		&s.Eruption,
		&s.Description,
		&s.EnglishDescription,
		&s.StandardAreaCode,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(categoryJSON), &s.Category)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func SaveShelter(s *model.Shelter) (*model.Shelter, error) {
	s.CreatedAt = time.Now()

	stmt, err := db.Prepare(QUERY_INSERT)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		s.FacilityName,
		s.EnglishFacilityName,
		s.CategoryJSON(),
		s.Prefecture,
		s.City,
		s.AdministrativeDistrict,
		s.Address,
		s.AddressKana,
		s.TelephoneNumber,
		s.TargetDistrict,
		s.TargetDistrictDescription,
		s.Latitude,
		s.Longitude,
		s.AboveSeaLevel,
		s.FacilityCapacity,
		s.FacilityHeight,
		s.Flood,
		s.Landslide,
		s.StormSurge,
		s.Earthquake,
		s.Tsunami,
		s.Inundation,
		s.Fire,
		s.Eruption,
		s.Description,
		s.EnglishDescription,
		s.StandardAreaCode,
		s.CreatedAt,
	).Scan(&s.ID)
	if err != nil {
		log.Println("failed to add shelter:", s.FacilityName)
		return nil, err
	}

	log.Println("added shelter:", s.FacilityName)

	return s, nil
}

func UpdateShelter(s *model.Shelter) (*model.Shelter, error) {
	ts := time.Now()
	s.UpdatedAt = &ts

	stmt, err := db.Prepare(QUERY_UPDATE)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	stmt.Exec(
		s.EnglishFacilityName,
		s.CategoryJSON(),
		s.Prefecture,
		s.City,
		s.AdministrativeDistrict,
		s.Address,
		s.AddressKana,
		s.TelephoneNumber,
		s.TargetDistrict,
		s.TargetDistrictDescription,
		s.Latitude,
		s.Longitude,
		s.AboveSeaLevel,
		s.FacilityCapacity,
		s.FacilityHeight,
		s.Flood,
		s.Landslide,
		s.StormSurge,
		s.Earthquake,
		s.Tsunami,
		s.Inundation,
		s.Fire,
		s.Eruption,
		s.Description,
		s.EnglishDescription,
		s.StandardAreaCode,
		s.UpdatedAt,
		s.FacilityName,
	)
	if err != nil {
		log.Println("failed to update shelter:", s.FacilityName)
		return nil, err
	}

	log.Println("updated shelter:", s.FacilityName)

	return s, nil
}

func DeleteShelter(s *model.Shelter) error {
	stmt, err := db.Prepare(QUERY_DELETE)
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(s.FacilityName)
	if err != nil {
		log.Println("failed to delete shelter:", s.FacilityName)
		return err
	}

	log.Println("deleted shelter:", s.FacilityName)

	return nil
}

func MergeShelter(s *model.Shelter) (*model.Shelter, error) {
	ex, err := FindShelter(s.FacilityName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("failed to merge shelter", err)
		return s, err
	}

	// add new shelter
	if ex == nil {
		SaveShelter(s)
		return s, nil
	}

	// update an existing shelter
	if s.EnglishFacilityName == nil {
		s.EnglishFacilityName = ex.EnglishFacilityName
	}
	for _, c := range ex.Category {
		s.AddCategory(c)
	}
	if s.Prefecture == nil {
		s.Prefecture = ex.Prefecture
	}
	if s.AdministrativeDistrict == nil {
		s.AdministrativeDistrict = ex.AdministrativeDistrict
	}
	if s.Address == nil {
		s.Address = ex.Address
	}
	if s.AddressKana == nil {
		s.AddressKana = ex.AddressKana
	}
	if s.TelephoneNumber == nil {
		s.TelephoneNumber = ex.TelephoneNumber
	}
	if s.TargetDistrict == nil {
		s.TargetDistrict = ex.TargetDistrict
	}
	if s.TargetDistrictDescription == nil {
		s.TargetDistrictDescription = ex.TargetDistrictDescription
	}
	if s.AboveSeaLevel == nil {
		s.AboveSeaLevel = ex.AboveSeaLevel
	}
	if s.FacilityCapacity == nil {
		s.FacilityCapacity = ex.FacilityCapacity
	}
	if s.FacilityHeight == nil {
		s.FacilityHeight = ex.FacilityHeight
	}
	if s.Flood == nil {
		s.Flood = ex.Flood
	}
	if s.Landslide == nil {
		s.Landslide = ex.Landslide
	}
	if s.StormSurge == nil {
		s.StormSurge = ex.StormSurge
	}
	if s.Earthquake == nil {
		s.Earthquake = ex.Earthquake
	}
	if s.Tsunami == nil {
		s.Tsunami = ex.Tsunami
	}
	if s.Inundation == nil {
		s.Inundation = ex.Inundation
	}
	if s.Fire == nil {
		s.Fire = ex.Fire
	}
	if s.Eruption == nil {
		s.Eruption = ex.Eruption
	}
	if s.Description == nil {
		s.Description = ex.Description
	}
	if s.EnglishDescription == nil {
		s.EnglishDescription = ex.EnglishDescription
	}
	if s.StandardAreaCode == nil {
		s.StandardAreaCode = ex.StandardAreaCode
	}
	UpdateShelter(s)

	return s, nil
}

package model

import (
	"strconv"
	"time"
)

const (
	facilityName = iota
	englishFacilityName
	category
	prefecture
	city
	administrativeDistrict
	address
	addressKana
	telephoneNumber
	targetDistrict
	targetDistrictDescription
	latitude
	longitude
	aboveSeaLevel
	facilityCapacity
	facilityHeight
	flood
	landslide
	stormSurge
	earthquake
	tsunami
	inundation
	fire
	eruption
	description
	englishDescription
	standardAreaCode
)

type Shelter struct {
	ID                        int       `json:"id"`
	FacilityName              string    `json:"facility_name"`
	EnglishFacilityName       string    `json:"facility_name_en"`
	Category                  string    `json:"category"`
	Prefecture                string    `json:"prefecture"`
	City                      string    `json:"city"`
	AdministrativeDistrict    string    `json:"administrative_district"`
	Address                   string    `json:"address"`
	AddressKana               string    `json:"address_kana"`
	TelephoneNumber           string    `json:"telephone_number"`
	TargetDistrict            string    `json:"target_district"`
	TargetDistrictDescription string    `json:"target_district_description"`
	Latitude                  float64   `json:"latitude"`
	Longitude                 float64   `json:"longitude"`
	AboveSeaLevel             float64   `json:"above_sea_level"`
	FacilityCapacity          int       `json:"facility_capacity"`
	FacilityHeight            float64   `json:"facility_height"`
	Flood                     bool      `json:"flood"`
	Landslide                 bool      `json:"landslide"`
	StormSurge                bool      `json:"storm_surge"`
	Earthquake                bool      `json:"earthquake"`
	Tsunami                   bool      `json:"tsunami"`
	Inundation                bool      `json:"inundation"`
	Fire                      bool      `json:"fire"`
	Eruption                  bool      `json:"eruption"`
	Description               string    `json:"description"`
	EnglishDescription        string    `json:"description_en"`
	StandardAreaCode          string    `json:"standard_area_code"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

func floats(s string, def float64) float64 {
	if len(s) == 0 {
		return def
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return def
	}

	return f
}

func ints(s string, def int) int {
	return int(floats(s, 0))
}

func bools(s string, def bool) bool {
	if len(s) == 0 {
		return def
	}
	return true
}

func NewShelter(record []string) *Shelter {
	return &Shelter{
		FacilityName:              record[facilityName],
		EnglishFacilityName:       record[englishFacilityName],
		Category:                  record[category],
		Prefecture:                record[prefecture],
		City:                      record[city],
		AdministrativeDistrict:    record[administrativeDistrict],
		Address:                   record[address],
		AddressKana:               record[addressKana],
		TelephoneNumber:           record[telephoneNumber],
		TargetDistrict:            record[targetDistrict],
		TargetDistrictDescription: record[targetDistrictDescription],
		Latitude:                  floats(record[latitude], 0),
		Longitude:                 floats(record[longitude], 0),
		AboveSeaLevel:             floats(record[aboveSeaLevel], 0),
		FacilityCapacity:          ints(record[facilityCapacity], 0),
		FacilityHeight:            floats(record[facilityHeight], 0),
		Flood:                     bools(record[flood], false),
		Landslide:                 bools(record[landslide], false),
		StormSurge:                bools(record[stormSurge], false),
		Earthquake:                bools(record[earthquake], false),
		Tsunami:                   bools(record[tsunami], false),
		Inundation:                bools(record[inundation], false),
		Fire:                      bools(record[fire], false),
		Eruption:                  bools(record[eruption], false),
		Description:               record[description],
		EnglishDescription:        record[englishDescription],
		StandardAreaCode:          record[standardAreaCode],
	}
}

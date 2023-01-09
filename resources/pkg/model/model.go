package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/satoshiyamamoto/shinagawa-shelter/pkg/csv"
)

type Shelter struct {
	ID                        int        `json:"id"`
	FacilityName              string     `json:"facility_name"`
	EnglishFacilityName       *string    `json:"facility_name_en,omitempty"`
	Category                  []string   `json:"category,omitempty"`
	Prefecture                *string    `json:"prefecture,omitempty"`
	City                      *string    `json:"city,omitempty"`
	AdministrativeDistrict    *string    `json:"administrative_district,omitempty"`
	Address                   *string    `json:"address,omitempty"`
	AddressKana               *string    `json:"address_kana,omitempty"`
	TelephoneNumber           *string    `json:"telephone_number,omitempty"`
	TargetDistrict            *string    `json:"target_district,omitempty"`
	TargetDistrictDescription *string    `json:"target_district_description,omitempty"`
	Latitude                  float64    `json:"latitude"`
	Longitude                 float64    `json:"longitude"`
	AboveSeaLevel             *float64   `json:"above_sea_level,omitempty"`
	FacilityCapacity          *int       `json:"facility_capacity,omitempty"`
	FacilityHeight            *float64   `json:"facility_height,omitempty"`
	Flood                     *bool      `json:"flood,omitempty"`
	Landslide                 *bool      `json:"landslide,omitempty"`
	StormSurge                *bool      `json:"storm_surge,omitempty"`
	Earthquake                *bool      `json:"earthquake,omitempty"`
	Tsunami                   *bool      `json:"tsunami,omitempty"`
	Inundation                *bool      `json:"inundation,omitempty"`
	Fire                      *bool      `json:"fire,omitempty"`
	Eruption                  *bool      `json:"eruption,omitempty"`
	Description               *string    `json:"description,omitempty"`
	EnglishDescription        *string    `json:"description_en,omitempty"`
	StandardAreaCode          *string    `json:"standard_area_code,omitempty"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 *time.Time `json:"updated_at,omitempty"`
}

func (s *Shelter) CategoryJSON() string {
	b, err := json.Marshal(s.Category)
	if err != nil {
		log.Println(err)
		return "[]"
	}
	return string(b)
}

func (s *Shelter) AddCategory(c string) bool {
	for _, v := range s.Category {
		if c == v {
			return false
		}
	}
	s.Category = append(s.Category, c)
	return true
}

func NewShelter(record []string) *Shelter {
	category := []string{
		*csv.Strings(record[csv.Category]),
	}

	return &Shelter{
		FacilityName:              *csv.Strings(record[csv.FacilityName]),
		EnglishFacilityName:       csv.Strings(record[csv.EnglishFacilityName]),
		Category:                  category,
		Prefecture:                csv.Strings(record[csv.Prefecture]),
		City:                      csv.Strings(record[csv.City]),
		AdministrativeDistrict:    csv.Strings(record[csv.AdministrativeDistrict]),
		Address:                   csv.Strings(record[csv.Address]),
		AddressKana:               csv.Strings(record[csv.AddressKana]),
		TelephoneNumber:           csv.Strings(record[csv.TelephoneNumber]),
		TargetDistrict:            csv.Strings(record[csv.TargetDistrict]),
		TargetDistrictDescription: csv.Strings(record[csv.TargetDistrictDescription]),
		Latitude:                  *csv.Floats(record[csv.Latitude]),
		Longitude:                 *csv.Floats(record[csv.Longitude]),
		AboveSeaLevel:             csv.Floats(record[csv.AboveSeaLevel]),
		FacilityCapacity:          csv.Ints(record[csv.FacilityCapacity]),
		FacilityHeight:            csv.Floats(record[csv.FacilityHeight]),
		Flood:                     csv.Bools(record[csv.Flood]),
		Landslide:                 csv.Bools(record[csv.Landslide]),
		StormSurge:                csv.Bools(record[csv.StormSurge]),
		Earthquake:                csv.Bools(record[csv.Earthquake]),
		Tsunami:                   csv.Bools(record[csv.Tsunami]),
		Inundation:                csv.Bools(record[csv.Inundation]),
		Fire:                      csv.Bools(record[csv.Fire]),
		Eruption:                  csv.Bools(record[csv.Eruption]),
		Description:               csv.Strings(record[csv.Description]),
		EnglishDescription:        csv.Strings(record[csv.EnglishDescription]),
		StandardAreaCode:          csv.Strings(record[csv.StandardAreaCode]),
	}
}

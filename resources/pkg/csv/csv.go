package csv

import (
	"strconv"
)

const (
	FacilityName = iota
	EnglishFacilityName
	Category
	Prefecture
	City
	AdministrativeDistrict
	Address
	AddressKana
	TelephoneNumber
	TargetDistrict
	TargetDistrictDescription
	Latitude
	Longitude
	AboveSeaLevel
	FacilityCapacity
	FacilityHeight
	Flood
	Landslide
	StormSurge
	Earthquake
	Tsunami
	Inundation
	Fire
	Eruption
	Description
	EnglishDescription
	StandardAreaCode
)

func Strings(s string) *string {
	if len(s) == 0 {
		return nil
	}
	return &s
}

func Floats(s string) *float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &f
}

func Ints(s string) *int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &i
}

func Bools(s string) *bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return nil
	}
	return &b
}

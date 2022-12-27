package model

type Shelter struct {
	ID                        int     `json:"id"`
	FacilityName              string  `json:"facility_name"`
	EnglishFacilityName       string  `json:"facility_name_en"`
	Category                  string  `json:"category"`
	Prefecture                string  `json:"prefecture"`
	City                      string  `json:"city"`
	AdministrativeDistrict    string  `json:"administrative_district"`
	Address                   string  `json:"address"`
	AddressKana               string  `json:"address_kana"`
	TelephoneNumber           string  `json:"telephone_number"`
	TargetDistrict            string  `json:"target_district"`
	TargetDistrictDescription string  `json:"target_district_description"`
	Latitude                  float64 `json:"latitude"`
	Longitude                 float64 `json:"longitude"`
	AboveSeaLevel             float64 `json:"above_sea_level"`
	FacilityCapacity          int     `json:"facility_capacity"`
	FacilityHeight            float64 `json:"facility_height"`
	Flood                     bool    `json:"flood"`
	Landslide                 bool    `json:"landslide"`
	StormSurge                bool    `json:"storm_surge"`
	Earthquake                bool    `json:"earthquake"`
	Tsunami                   bool    `json:"tsunami"`
	Inundation                bool    `json:"inundation"`
	Fire                      bool    `json:"fire"`
	Eruption                  bool    `json:"eruption"`
	Description               string  `json:"description"`
	EnglishDescription        string  `json:"description_en"`
	StandardAreaCode          string  `json:"standard_area_code"`
}

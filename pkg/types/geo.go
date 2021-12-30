package types

type GeoIP struct {
	IP          string  `json:"ip,omitempty"`
	CountryCode string  `json:"country_code,omitempty"`
	CountryName string  `json:"country_name,omitempty"`
	RegionCode  string  `json:"region_code,omitempty"`
	RegionName  string  `json:"region_name,omitempty"`
	City        string  `json:"city,omitempty"`
	Zipcode     string  `json:"zip_code,omitempty"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
	MetroCode   int64   `json:"metro_code,omitempty"`
	AreaCode    int64   `json:"area_code,omitempty"`
	Active      int64   `json:"active,omitempty"`
}

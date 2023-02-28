package models

// {"country_code":"US","country_name":"United States","city":"Los Angeles","postal":"90007","latitude":34.0294,"longitude":-118.2871,"IPv4":"207.151.52.114","state":"California"}

type Location struct {
	CountryCode string  `json:"country_code" bson:"countryCode,omitempty"`
	CountryName string  `json:"country_name" bson:"countryName,omitempty"`
	City        string  `json:"city" bson:"city,omitempty"`
	Postal      string  `json:"postal" bson:"postal,omitempty"`
	Latitude    float64 `json:"latitude" bson:"latitude,omitempty"`
	Longitude   float64 `json:"longitude" bson:"longitude,omitempty"`
	State       string  `json:"state" bson:"state,omitempty"`
}

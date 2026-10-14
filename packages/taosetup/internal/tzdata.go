package internal

import (
	_ "embed"
	"encoding/json"
)

type timezoneCountry struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type timezoneCity struct {
	Name    string          `json:"name"`
	Country timezoneCountry `json:"country"`
}

//go:embed tzdata.json
var tzdata []byte
var timezoneOptions map[string][]timezoneCity

func GetTimezoneOptions() map[string][]timezoneCity {
	if len(timezoneOptions) == 0 {
		timezoneOptions = make(map[string][]timezoneCity)
		err := json.Unmarshal(tzdata, &timezoneOptions)
		if err != nil {
			panic(err)
		}
	}
	return timezoneOptions
}

func GetCountryCode(region, city string) string {
	for _, entry := range timezoneOptions[region] {
		if entry.Name == city {
			return entry.Country.Code
		}
	}
	return ""
}

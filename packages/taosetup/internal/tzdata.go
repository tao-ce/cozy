package internal

import (
	_ "embed"
	"encoding/json"
)

//go:embed tzdata.json
var tzdata []byte
var timezoneOptions map[string][]string

func GetTimezoneOptions() map[string][]string {
	if len(timezoneOptions) == 0 {
		timezoneOptions = make(map[string][]string)
		err := json.Unmarshal(tzdata, &timezoneOptions)
		if err != nil {
			panic(err)
		}
	}
	return timezoneOptions
}

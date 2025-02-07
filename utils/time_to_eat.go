package utils

import (
	"fmt"
	"time"
)

// ConvertToEAT converts a UTC time string to East African Time (EAT)
func ConvertToEAT(utcTime string) (time.Time, error) {
	// parse the time (assuming UTC in RFC3339 format)
	parsedTime, err := time.Parse(time.RFC3339, utcTime)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing time: %v", err)
	}

	// load the EAT timezone (Africa/Nairobi)
	eatLocation, err := time.LoadLocation("Africa/Nairobi")
	if err != nil {
		return time.Time{}, fmt.Errorf("error loading EAT timezone: %v", err)
	}

	// convert the UTC time to EAT
	eatTime := parsedTime.In(eatLocation)

	return eatTime, nil
}

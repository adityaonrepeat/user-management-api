package service

import "time"

func CalculateAge(dob, now time.Time) int {
	years := now.Year() - dob.Year()

	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		years--
	}

	return years
}

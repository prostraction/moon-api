package server

import (
	"errors"
	"math"
	"strconv"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func StrToInt(val string, fallback int, min int, max int) int {
	v, err := strconv.Atoi(val)
	if err != nil {
		v = fallback
	}
	if v > max {
		v = max
	}
	if v < min {
		v = min
	}

	return v
}

func parseCoords(latStr, lonStr string) Coordinates {
	locationCords := Coordinates{}
	if latStr == "no-value" || lonStr == "no-value" {
		locationCords.IsValid = false
	} else {
		lat, err1 := strconv.ParseFloat(latStr, 64)
		lon, err2 := strconv.ParseFloat(lonStr, 64)
		if err1 != nil || err2 != nil {
			locationCords.IsValid = false
		} else {
			locationCords.IsValid = true
			locationCords.Latitude = lat
			locationCords.Longitude = lon
		}
	}
	return locationCords
}

func IsValidDate(year, month, day int) error {
	if year < 0 || year > 9999 {
		return errors.New("'year' should be in range [0,9999]")
	}

	if month < 1 || month > 12 {
		return errors.New("'month' should be in range [1,12]")
	}

	if day < 1 {
		return errors.New("'day' should be greater then 0")
	}

	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	if month == 2 {
		if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
			daysInMonth[1] = 29
		} else {
			daysInMonth[1] = 28
		}
	}

	if day > daysInMonth[month-1] {
		return errors.New("no 'day' in 'month' this 'year'")
	}

	return nil
}

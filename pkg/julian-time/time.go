package julian_time

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	m "moon/pkg/math-helpers"
)

const (
	HoursPerDay      = 24.0
	MinutesPerHour   = 60.0
	SecondsPerMinute = 60.0
)

func ToJulianDate(t time.Time) float64 {
	year := t.Year()
	month := int(t.Month())

	fractionalDay := (float64(t.Hour()) +
		float64(t.Minute())/MinutesPerHour +
		float64(t.Second())/(MinutesPerHour*SecondsPerMinute)) / HoursPerDay

	day := float64(t.Day()) + fractionalDay

	if month <= 2 {
		year -= 1
		month += 12
	}

	a := year / 100
	b := 2 - a + (a / 4)

	jd := math.Floor(365.25*float64(year+4716)) +
		math.Floor(30.6001*float64(month+1)) +
		day + float64(b) - 1524.5

	return jd
}

func FromJulianDate(j float64, loc *time.Location) time.Time {
	datey, datem, dated := Jyear(j)
	timeh, timem, times := Jhms(j)

	t := time.Date(datey, GetMonth(datem), dated, timeh, timem, times, 0, time.UTC)
	t = t.In(loc)
	return t
}

func SetTimezoneLocFromString(utc string) (*time.Location, error) {
	if utc == "" {
		return time.UTC, nil
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9:+-\-]`)
	utc = re.ReplaceAllString(utc, "")

	normalized := strings.ToLower(utc)
	normalized = strings.TrimPrefix(normalized, "utc")
	normalized = strings.TrimPrefix(normalized, "gmt")
	normalized = strings.TrimSpace(normalized)

	if normalized == "" || normalized == "+0" || normalized == "-0" || normalized == "0" {
		return time.UTC, nil
	}

	sign := 1
	if strings.HasPrefix(normalized, "+") {
		sign = 1
		normalized = normalized[1:]
	} else if strings.HasPrefix(normalized, "-") {
		sign = -1
		normalized = normalized[1:]
	}

	var hours, minutes int
	var err error

	if strings.Contains(normalized, ":") {
		parts := strings.Split(normalized, ":")
		if len(parts) != 2 {
			return time.UTC, fmt.Errorf("invalid timezone format: %s", utc)
		}

		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			return time.UTC, fmt.Errorf("invalid hours: %s", parts[0])
		}

		minutes, err = strconv.Atoi(parts[1])
		if err != nil || minutes < 0 || minutes >= 60 {
			return time.UTC, fmt.Errorf("invalid minutes: %s", parts[1])
		}

	} else {
		switch len(normalized) {
		case 1, 2:
			hours, err = strconv.Atoi(normalized)
			if err != nil {
				return time.UTC, fmt.Errorf("invalid hours: %s", normalized)
			}
			minutes = 0

		case 3:
			hours, err = strconv.Atoi(normalized[:1])
			if err != nil {
				return time.UTC, fmt.Errorf("invalid hours: %s", normalized[:1])
			}
			minutes, err = strconv.Atoi(normalized[1:])
			if err != nil || minutes < 0 || minutes >= 60 {
				return time.UTC, fmt.Errorf("invalid minutes: %s", normalized[1:])
			}

		case 4:
			hours, err = strconv.Atoi(normalized[:2])
			if err != nil {
				return time.UTC, fmt.Errorf("invalid hours: %s", normalized[:2])
			}
			minutes, err = strconv.Atoi(normalized[2:])
			if err != nil || minutes < 0 || minutes >= 60 {
				return time.UTC, fmt.Errorf("invalid minutes: %s", normalized[2:])
			}

		default:
			return time.UTC, fmt.Errorf("invalid timezone format: %s", utc)
		}
	}

	if hours < 0 || hours > 23 {
		return time.UTC, fmt.Errorf("hours out of range (0-23): %d", hours)
	}

	totalSeconds := sign * (hours*3600 + minutes*60)

	locationName := fmt.Sprintf("UTC%s%d:%02d", m.GetSignPrefix(sign), hours, minutes)
	if minutes == 0 {
		locationName = fmt.Sprintf("UTC%s%d", m.GetSignPrefix(sign), hours)
	}

	return time.FixedZone(locationName, totalSeconds), nil
}

func GetTimeFromLocation(loc *time.Location) (hours int, minutes int, err error) {
	if loc == nil {
		return 0, 0, errors.New("loc is nil")
	}
	utc := loc.String()
	re := regexp.MustCompile(`[^a-zA-Z0-9:+-\-]`)
	utc = re.ReplaceAllString(utc, "")

	if utc == "" {
		return 0, 0, errors.New("empty timezone string")
	}
	normalized := strings.ToLower(utc)
	normalized = strings.TrimPrefix(normalized, "utc")
	normalized = strings.TrimPrefix(normalized, "gmt")
	normalized = strings.TrimSpace(normalized)

	if normalized == "" || normalized == "+0" || normalized == "-0" || normalized == "0" {
		return 0, 0, nil
	}

	var sign int = 1
	if strings.HasPrefix(normalized, "+") {
		sign = 1
		normalized = normalized[1:]
	} else if strings.HasPrefix(normalized, "-") {
		sign = -1
		normalized = normalized[1:]
	}

	if strings.Contains(normalized, ":") {
		parts := strings.Split(normalized, ":")
		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("invalid timezone format: %s", utc)
		}

		hours, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid hours: %s", parts[0])
		}

		minutes, err := strconv.Atoi(parts[1])
		if err != nil || minutes < 0 || minutes >= 60 {
			return 0, 0, fmt.Errorf("invalid minutes: %s", parts[1])
		}

		return sign * hours, minutes, nil
	}

	if len(normalized) <= 2 {
		hours, err := strconv.Atoi(normalized)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid hours: %s", normalized)
		}
		return sign * hours, minutes, nil
	}

	if len(normalized) == 4 {
		hoursStr := normalized[:2]
		minutesStr := normalized[2:]

		hours, err := strconv.Atoi(hoursStr)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid hours: %s", hoursStr)
		}

		minutes, err := strconv.Atoi(minutesStr)
		if err != nil || minutes < 0 || minutes >= 60 {
			return 0, 0, fmt.Errorf("invalid minutes: %s", minutesStr)
		}

		return sign * hours, minutes, nil
	}

	if len(normalized) == 3 {
		hoursStr := normalized[:1]
		minutesStr := normalized[1:]

		hours, err := strconv.Atoi(hoursStr)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid hours: %s", hoursStr)
		}

		minutes, err := strconv.Atoi(minutesStr)
		if err != nil || minutes < 0 || minutes >= 60 {
			return 0, 0, fmt.Errorf("invalid minutes: %s", minutesStr)
		}

		return sign * hours, minutes, nil
	}

	return 0, 0, fmt.Errorf("invalid timezone format")
}

// JYMD - Convert Julian time to year, months, and days
func Jyear(td float64) (int, int, int) {
	td += 0.5 // Astronomical to civil
	z := math.Floor(td)
	f := td - z

	var a float64
	if z < 2299161.0 {
		a = z
	} else {
		alpha := math.Floor((z - 1867216.25) / 36524.25)
		a = z + 1 + alpha - math.Floor(alpha/4)
	}

	b := a + 1524
	c := math.Floor((b - 122.1) / 365.25)
	d := math.Floor(365.25 * c)
	e := math.Floor((b - d) / 30.6001)

	mm := int(math.Floor(e))
	if mm >= 14 {
		mm -= 13
	} else {
		mm -= 1
	}

	year := int(math.Floor(c))
	if mm > 2 {
		year -= 4716
	} else {
		year -= 4715
	}

	day := int(math.Floor(b - d - math.Floor(30.6001*e) + f))

	return year, mm, day
}

// JHMS - Convert Julian time to hour, minutes, and seconds
func Jhms(j float64) (int, int, int) {
	j += 0.5 // Astronomical to civil
	ij := (j - math.Floor(j)) * 86400.0
	hours := math.Floor(ij / 3600)
	minutes := math.Floor((ij / 60))
	seconds := math.Floor(ij)
	return int(hours), int(math.Mod(minutes, 60)), int(math.Mod(seconds, 60))
}

func GetMonth(datem int) time.Month {
	datem = min(max(datem-1, 0), 11)
	return months[datem]
}

package position

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	jt "moon/pkg/julian-time"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Cache struct {
	CacheDaily   map[string]*DayData
	CacheMonthly map[string]*[]DayData
}

// DayResponse
type DayResponse struct {
	Status     string     `json:"Status"`
	Parameters Parameters `json:"Parameters"`
	Data       *DayData   `json:"Data"`
	Range      string     `json:"Range"`
}

// MonthResponse
type MonthResponse struct {
	Status     string     `json:"Status"`
	Message    string     `json:"Message,omitempty"`
	Parameters Parameters `json:"Parameters"`
	Data       []DayData  `json:"Data"`
	Range      string     `json:"Range"`
	DaysCount  int        `json:"DaysCount"`
}

// input
type Parameters struct {
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
	Timezone  int     `json:"Timezone"`
	UTCOffset string  `json:"UtcOffset"`
	Year      int     `json:"Year"`
	Month     int     `json:"Month"`
	Day       int     `json:"Day,omitempty"`
}

// resp for 1 day
type MoonPosition struct {
	Timestamp       *int64  `json:"Timestamp,omitempty"`
	Time            *any    `json:"Time"`
	AzimuthDegrees  float64 `json:"AzimuthDegrees"`
	AltitudeDegrees float64 `json:"AltitudeDegrees"`
	Direction       string  `json:"Direction"`
	DistanceKm      float64 `json:"DistanceKm"`
}

type DayData struct {
	Day        *string       `json:"Date,omitempty"`
	IsMoonRise bool          `json:"IsMoonRise"`
	IsMoonSet  bool          `json:"IsMoonSet"`
	IsMeridian bool          `json:"IsMeridian"`
	Moonrise   *MoonPosition `json:"Moonrise,omitempty"`
	Moonset    *MoonPosition `json:"Moonset,omitempty"`
	Meridian   *MoonPosition `json:"Meridian,omitempty"`
}

func (c *Cache) GetRisesMonthly(year, month int, loc *time.Location, precision int, timeFormat string, location ...float64) (*[]DayData, error) {
	lat, lon, err := parseLocation(location)
	if err != nil {
		return nil, err
	}

	if c.CacheMonthly == nil {
		c.CacheMonthly = make(map[string]*[]DayData)
	}

	var strKey strings.Builder
	strKey.WriteString(strconv.Itoa(year))
	strKey.WriteString("-")
	strKey.WriteString(strconv.Itoa(month))
	strKey.WriteString("-")
	if loc != nil {
		strKey.WriteString(loc.String())
	} else {
		strKey.WriteString("nil")
	}
	strKey.WriteString("-")
	strKey.WriteString(strconv.Itoa(precision))
	strKey.WriteString("-")
	strKey.WriteString(strconv.FormatFloat(lat, 'e', precision, 64))
	strKey.WriteString("-")
	strKey.WriteString(strconv.FormatFloat(lon, 'e', precision, 64))
	strKey.WriteString("-")
	strKey.WriteString(timeFormat)

	if c.CacheMonthly != nil && c.CacheMonthly[strKey.String()] != nil {
		return c.CacheMonthly[strKey.String()], nil
	}

	h := 0
	if loc != nil {
		jth, _, err := jt.GetTimeFromLocation(loc)
		if err == nil {
			h = jth
		}
	}

	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%.2f", lat))
	params.Add("lon", fmt.Sprintf("%.2f", lon))
	params.Add("utc", fmt.Sprintf("%d", h))
	params.Add("year", fmt.Sprintf("%d", year))
	params.Add("month", fmt.Sprintf("%d", month))
	params.Add("precision", fmt.Sprintf("%d", precision))

	url := baseURL + "monthly?" + params.Encode()
	client := &http.Client{Timeout: 69 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("[%s] Failed to make request: %w", resp.Status, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[%s] Failed to read response: %w", resp.Status, err)
	}

	var monthResponse MonthResponse
	if err := json.Unmarshal(body, &monthResponse); err != nil {
		return nil, fmt.Errorf("[%s] Failed to unmarshal response: %w", resp.Status, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[%s] %s", resp.Status, monthResponse.Message)
	}

	for i := range monthResponse.Data {
		if monthResponse.Data[i].Meridian != nil && monthResponse.Data[i].Meridian.Timestamp != nil {
			var t any = timestampToGoTime(monthResponse.Data[i].Meridian.Timestamp, timeFormat, loc)
			monthResponse.Data[i].Meridian.Time = &t
			monthResponse.Data[i].Meridian.Timestamp = nil
		}
		if monthResponse.Data[i].Moonrise != nil && monthResponse.Data[i].Moonrise.Timestamp != nil {
			var t any = timestampToGoTime(monthResponse.Data[i].Moonrise.Timestamp, timeFormat, loc)
			monthResponse.Data[i].Moonrise.Time = &t
			monthResponse.Data[i].Moonrise.Timestamp = nil
		}
		if monthResponse.Data[i].Moonset != nil {
			var t any = timestampToGoTime(monthResponse.Data[i].Moonset.Timestamp, timeFormat, loc)
			monthResponse.Data[i].Moonset.Time = &t
			monthResponse.Data[i].Moonset.Timestamp = nil
		}
		t := time.Date(year, jt.GetMonth(month), 1+i, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		monthResponse.Data[i].Day = &t
	}

	if c.CacheMonthly != nil && c.CacheMonthly[strKey.String()] == nil {
		c.CacheMonthly[strKey.String()] = &monthResponse.Data
	}

	return &monthResponse.Data, nil
}

func GetRisesDay(year, month, day int, loc *time.Location, precision int, timeFormat string, location ...float64) (*DayData, error) {
	lat, lon, err := parseLocation(location)
	if err != nil {
		return nil, err
	}

	h, m := 0, 0
	if loc != nil {
		jth, jtm, err := jt.GetTimeFromLocation(loc)
		if err == nil {
			h = jth
			m = jtm
		}
	}

	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%.2f", lat))
	params.Add("lon", fmt.Sprintf("%.2f", lon))
	params.Add("utc_hours", fmt.Sprintf("%d", h))
	params.Add("utc_minutes", fmt.Sprintf("%d", m))
	params.Add("year", fmt.Sprintf("%d", year))
	params.Add("month", fmt.Sprintf("%d", month))
	params.Add("day", fmt.Sprintf("%d", day))
	params.Add("precision", fmt.Sprintf("%d", precision))

	url := baseURL + "daily" + "?" + params.Encode()
	client := &http.Client{Timeout: 69 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var dayResponse DayResponse
	if err := json.Unmarshal(body, &dayResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if dayResponse.Data.Meridian != nil && dayResponse.Data.Meridian.Timestamp != nil {
		var t any = timestampToGoTime(dayResponse.Data.Meridian.Timestamp, timeFormat, loc)
		dayResponse.Data.Meridian.Time = &t
		dayResponse.Data.Meridian.Timestamp = nil
	}
	if dayResponse.Data.Moonrise != nil && dayResponse.Data.Moonrise.Timestamp != nil {
		var t any = timestampToGoTime(dayResponse.Data.Moonrise.Timestamp, timeFormat, loc)
		dayResponse.Data.Moonrise.Time = &t
		dayResponse.Data.Moonrise.Timestamp = nil
	}
	if dayResponse.Data.Moonset != nil {
		var t any = timestampToGoTime(dayResponse.Data.Moonset.Timestamp, timeFormat, loc)
		dayResponse.Data.Moonset.Time = &t
		dayResponse.Data.Moonset.Timestamp = nil
	}

	return dayResponse.Data, nil
}

func (c *Cache) GetRisesDay(year, month, day int, loc *time.Location, precision int, timeFormat string, location ...float64) (*DayData, error) {
	lat, lon, err := parseLocation(location)
	if err != nil {
		return nil, err
	}

	if c.CacheDaily == nil {
		c.CacheDaily = make(map[string]*DayData)
	}

	var strKey strings.Builder
	strKey.WriteString(strconv.Itoa(year))
	strKey.WriteString("-")
	strKey.WriteString(strconv.Itoa(month))
	strKey.WriteString("-")
	strKey.WriteString(strconv.Itoa(day))
	strKey.WriteString("-")
	if loc != nil {
		strKey.WriteString(loc.String())
	} else {
		strKey.WriteString("nil")
	}
	strKey.WriteString("-")
	strKey.WriteString(strconv.Itoa(precision))
	strKey.WriteString("-")
	strKey.WriteString(strconv.FormatFloat(lat, 'e', precision, 64))
	strKey.WriteString("-")
	strKey.WriteString(strconv.FormatFloat(lon, 'e', precision, 64))
	strKey.WriteString("-")
	strKey.WriteString(timeFormat)

	if c.CacheDaily != nil && c.CacheDaily[strKey.String()] != nil {
		return c.CacheDaily[strKey.String()], nil
	}

	dayResponse, err := GetRisesDay(year, month, day, loc, precision, timeFormat, location...)
	if err != nil {
		return nil, err
	}

	if c.CacheDaily != nil && c.CacheDaily[strKey.String()] == nil {
		c.CacheDaily[strKey.String()] = dayResponse
	}

	return dayResponse, nil
}

func GetMoonPosition(tGiven time.Time, loc *time.Location, precision int, timeFormat string, location ...float64) (*MoonPosition, error) {
	lat, lon, err := parseLocation(location)
	if err != nil {
		return nil, err
	}

	h, m := 0, 0
	if loc != nil {
		tGiven = tGiven.In(loc)
		jth, jtm, err := jt.GetTimeFromLocation(loc)
		if err == nil {
			h = jth
			m = jtm
		}
	}

	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%.2f", lat))
	params.Add("lon", fmt.Sprintf("%.2f", lon))
	params.Add("utc_hours", fmt.Sprintf("%d", h))
	params.Add("utc_minutes", fmt.Sprintf("%d", m))
	params.Add("year", fmt.Sprintf("%d", tGiven.Year()))
	params.Add("month", fmt.Sprintf("%d", int(tGiven.Month())))
	params.Add("day", fmt.Sprintf("%d", tGiven.Day()))
	params.Add("hour", fmt.Sprintf("%d", tGiven.Hour()))
	params.Add("minute", fmt.Sprintf("%d", tGiven.Minute()))
	params.Add("second", fmt.Sprintf("%d", tGiven.Second()))
	params.Add("precision", fmt.Sprintf("%d", precision))

	url := baseURL + "position" + "?" + params.Encode()
	client := &http.Client{Timeout: 69 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status: %s (%s)", resp.Status, resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var pos *MoonPosition
	if err := json.Unmarshal(body, &pos); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if pos != nil && pos.Timestamp != nil {
		var t any = timestampToGoTime(pos.Timestamp, timeFormat, loc)
		pos.Time = &t
		pos.Timestamp = nil
	}

	return pos, nil
}

func parseLocation(location []float64) (lat, lon float64, err error) {
	if len(location) == 2 {
		lat = location[1]
		lon = location[0]
	} else {
		return 0, 0, errors.New("no location prodived")
	}

	return lat, lon, nil
}

func timestampToGoTime(ev *int64, timeFormat string, loc *time.Location) *any {
	if ev == nil {
		return nil
	}
	utcTime := time.Unix(*ev, 0).UTC()
	time := utcTime
	if loc != nil {
		time = time.In(loc)
	}

	if strings.ToLower(timeFormat) == "timestamp" {
		var t any = time.Unix()
		return &t
	}

	if strings.ToLower(timeFormat) != "iso" {
		var t any = time.Format(timeFormat)
		return &t
	}

	var t any = time
	return &t
}

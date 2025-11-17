# API for Moon calculations

## Run

Port by default: 9998

Main server:
```
git clone https://github.com/prostraction/moon
cd moon
go run cmd/main.go
```

Addon:
```
pip install skyfield
python addon/skyfield/server.py
```

## Methods

### GET /api/v1/moonPhaseDate

The method returns the Moon parameters for the specified day and time. If the day or time is not specified, the current value for the unspecified fields is taken. If longitude and latitude are specified, the response will contain additional structures.

#### Params

  | Parameter | Type | Description | Example Value |
| :--- | :--- | :--- |  :--- | 
|`utc` | `string [optional, default="UTC+0"]` | UTC in format `UTC+7`, `UTC+09:30`, `-3` | `UTC+4`
|`lang` | `string [optional, default="en"]` | Values available: ("en", "es", "fr", "de", "ru", "jp") | `es`
|`precision` | `int [optional, default=2]` | How many digits after ```.``` will be in output. Allowed range: [1, 20] | `5`
|`timeFormat` | `string [optional]` | How the data will be formed. Available values: ISO (default), timestamp, duration or Go format. Duration is seconds between current date and event time.| `duration`
|`latitude` | `float [optional, default=none]` | Latitude of viewer's place. Used for moon position calculations: ```MoonDaysDetailed```, ```MoonRiseAndSet```, and ```MoonPosition``` object | `51.1655`
|`longitude` | `float [optional, default=none]` | Longitude of viewer's place. Used for moon position calculations: ```MoonDaysDetailed```, ```MoonRiseAndSet```, and ```MoonPosition``` object | `71.4272`
|`year` | `int [optional, default=<current year>]` | Format: YYYY Allowed range: [1, 9999] | `2025`
|`month` | `int [optional, default=<current month>]` | Format: M or MM. Allowed range: [1, 12] | `01` or `1`
|`day` | `int [optional, default=<current day>]` | Format: D or DD. Allowed range: [1, 31] | `01` or `1`
|`hour` | `int [optional, default=<current hour>]` | Format: h or hh. Allowed range: [0, 23] | `01` or `1`
|`minute` | `int [optional, default=<current minute>]` | Format: m or mm. Allowed range: [0, 59] | `01` or `1`
|`second` | `int [optional, default=<current second>]` | Format: s or ss. Allowed range: [0, 59] | `01` or `1`


----------------------------------------------------------------

#### Response:

The method returns 6 objects:
- ```BeginDay```, ```CurrentState```, ```EndDay``` objects of the ```MoonStat``` structure to display the position of the moon at the beginning of the day, the specified time and the end of the day, respectively;
- ```MoonDaysDetailed```, a structure for determining the number of lunar days on a given day;
- ```ZodiacDetailed```, a structure for determining which zodiac sign the moon is in on a given time interval, when it began and ended. It contains an array for each lunar day that falls on a given Earth day;
- ```MoonRiseAndSet```, a structure for determining the moonrise, moonset and meridian on a given day.

<details>
  <summary><strong>Table</strong></summary>
  
  | Response Variable | Type | Description |
| :--- | :--- | :--- | 
|`BeginDay` | `Object of struct MoonStat [required]` | Data for the beginning of the requested day (00:00) |
|`CurrentState` | `Object of struct MoonStat [required]` | Data at specified time of requested day: hour, minute and second from request Params. |
|`EndDay` | `Object of struct MoonStat [required]` | Data for end of requested day (00:00 next day) |
|`MoonDaysDetailed` | `Object of struct MoonDaysDetailed [optional]` | Detailed lunar day information that falls on a given Earth day. Exists only if latitude and longitude are specified. |
|`ZodiacDetailed` | `Object of struct ZodiacDetailed [required]` | Detailed zodiac transit information | 
|`MoonRiseAndSet` | `Object of struct MoonRiseAndSet [optional]` | Moon rise/set/meridian events. Exists only if latitude and longitude are specified. |
  
</details>

----------------------------------------------------------------

##### MoonStat (used as ```BeginDay```, ```CurrentState```, ```EndDay```)

```MoonStat``` objects are used to display at a given time.
In case of a method response, MoonStat will contain the values:
- BeginDay: start of the day, 00:00AM
- CurrentState: given time from request
- EndDay: start of the next day, 00:00AM

<details>
  <summary><strong>Table</strong></summary>

  | Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`MoonStat.MoonDays` | `Float [required]` | Lunar day number | `23.54` |
|`MoonStat.Illumination` | `Float [required]` | Percentage of Moon's disk illuminated | `38.27` |
|`MoonStat.Phase` | `Object of struct Phase [required]` | Lunar phase details | - |
|`MoonStat.Zodiac` | `Object of struct Zodiac [required]` | Zodiac sign details | - |
|`MoonStat.MoonPosition` | `Object of struct MoonPosition [optional]` | Moon position data. Exists only if latitude and longitude are specified. | - |

Phase structure:

  | Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`Phase.Name` | `String [required]` | Phase name in English | `"Waning Crescent"` |
|`Phase.NameLocalized` | `String [required]` | Localized phase name | `"Убывающий серп"` |
|`Phase.Emoji` | `String [required]` | Phase emoji | `"🌘"` |
|`Phase.IsWaxing` | `Boolean [required]` | True if Moon is waxing / illumination is increasing | `false` |

Zodiac structure:

  | Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`Zodiac.Name` | `String [required]` | Zodiac name in English | `"Gemini"` |
|`Zodiac.NameLocalized` | `String [required]` | Localized zodiac name | `"Близнецы"` |
|`Zodiac.Emoji` | `String [required]` | Zodiac emoji | `"♊"` |

MoonPosition structure (Exists only if latitude and longitude are specified):

  | Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`MoonPosition.Time` | `String / int64 [required]` | Time in ISO (default) or specified timeFormat | `"2025-09-16T00:00:00+05:00"` |
|`MoonPosition.AzimuthDegrees` | `Float [required]` | Compass direction (0°=North) | `57.1` |
|`MoonPosition.AltitudeDegrees` | `Float [required]` | Angle above horizon (negative = below) | `8.8` |
|`MoonPosition.Direction` | `String [required]` | Cardinal direction abbreviation | `"ENE"` |
|`MoonPosition.DistanceKm` | `Float [required]` | Earth-Moon distance in km | `376559.9` |

</details>

</details>

----------------------------------------------------------------

##### MoonDaysDetailed

```MoonDaysDetailed``` is a structure that contains an array for each lunar day that falls on a given Earth day. Exists only if latitude and longitude are specified.

<details>
  <summary><strong>Table</strong></summary>

  | Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`MoonDaysDetailed.Count` | `Integer [required]` | Number of lunar days this calendar day | `2` |
|`MoonDaysDetailed.Day` | `Array<Object> [required]` | Array of lunar day periods | - |
|`MoonDaysDetailed.Day[].Begin` | `String [optional]` | Start time of lunar day (ISO 8601) | `"2025-09-15T22:37:45+05:00"` |
|`MoonDaysDetailed.Day[].IsBeginExists` | `Boolean [required]` | True if start time is past/present | `true` |
|`MoonDaysDetailed.Day[].End` | `String [optional]` | End time of lunar day (ISO 8601) | `"2025-09-16T23:56:10+05:00"` |
|`MoonDaysDetailed.Day[].IsEndExists` | `Boolean [required]` | True if end time is past | `true`, `false` |

</details>

----------------------------------------------------------------

##### ZodiacDetailed

```ZodiacDetailed``` is a structure for determining which zodiac sign the moon is in on a given time interval, when it began and ended. It contains an array for each lunar day that falls on a given Earth day.

<details>
  <summary><strong>Table</strong></summary>

  | Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`ZodiacDetailed.Count` | `Integer [required]` | Number of zodiac signs this day | `1` |
|`ZodiacDetailed.Zodiac` | `Array<Object> [required]` | Array of zodiac transit periods | - |
|`ZodiacDetailed.Zodiac[].Name` | `String [required]` | Zodiac sign name | `"Gemini"` |
|`ZodiacDetailed.Zodiac[].NameLocalized` | `String [required]` | Localized zodiac name | `"Близнецы"` |
|`ZodiacDetailed.Zodiac[].Emoji` | `String [required]` | Zodiac emoji | `"♊"` |
|`ZodiacDetailed.Zodiac[].Begin` | `String [required]` | Entry time into sign (ISO 8601) | `"2025-09-14T23:07:06+05:00"` |
|`ZodiacDetailed.Zodiac[].End` | `String [required]` | Exit time from sign (ISO 8601) | `"2025-09-17T11:07:06+05:00"` |

</details>

----------------------------------------------------------------

##### MoonRiseAndSet

```MoonRiseAndSet``` is a structure for determining the moonrise, moonset and meridian on a given day. Exists only if latitude and longitude are specified.

<details>
  <summary><strong>Table</strong></summary>

| Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`MoonRiseAndSet.Date` | `String [optional]` | Date for day of calculations, missing as default | `2025-01-15` |
|`MoonRiseAndSet.IsMoonRise` | `Boolean [required]` | True if moonrise occurs at given day | `true` |
|`MoonRiseAndSet.IsMoonSet` | `Boolean [required]` | True if moonset occurs at given day | `true` |
|`MoonRiseAndSet.IsMeridian` | `Boolean [required]` | True if meridian transit occurs at given day | `true` |
|`MoonRiseAndSet.Moonrise` | `Object of struct MoonPosition [optional]` | Moonrise position data. Exists only if IsMoonRise = true | - |
|`MoonRiseAndSet.Moonset` | `Object of struct MoonPosition [optional]` | Moonset position data. Exists only if IsMoonSet = true | - |
|`MoonRiseAndSet.Meridian` | `Object of struct MoonPosition [optional]` | Meridian position data, Exists only if IsMeridian = true | - |

MoonPosition structure:

<details>
  <summary><strong>Table</strong></summary>

| Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`MoonPosition.Time` | `String / int64 [required]` | Time in ISO (default) or specified timeFormat | `"2025-09-16T00:00:00+05:00"` |
|`MoonPosition.AzimuthDegrees` | `Float [required]` | Moonrise azimuth | `47.3` |
|`MoonPosition.AltitudeDegrees` | `Float [required]` | Moonrise altitude | `-0.6` |
|`MoonPosition.Direction` | `String [required]` | Moonrise direction | `"ENE"` |
|`MoonPosition.DistanceKm` | `Float [required]` | Earth-Moon distance in km | `376559.9` |

</details>

</details>

----------------------------------------------------------------

#### Response examples

Response of:
```GET /api/v1/moonPhaseDate?lang=ru&utc=5&latitude=51.1655&longitude=71.4272&year=2025&month=09&day=15&precision=5&hour=12&minute=0&second=0```

<details>
  <summary><strong>JSON</strong></summary>

```json
{
  "BeginDay": {
    "MoonDays": 22.53674,
    "Illumination": 49.42951,
    "Phase": {
      "Name": "Third quarter",
      "NameLocalized": "Последняя четверть",
      "Emoji": "🌗",
      "IsWaxing": false
    },
    "Zodiac": {
      "Name": "Gemini",
      "NameLocalized": "Близнецы",
      "Emoji": "♊"
    },
    "MoonPosition": {
      "Time": "2025-09-15T02:00:00+05:00",
      "AzimuthDegrees": 66.96877,
      "AltitudeDegrees": 17.50543,
      "Direction": "ENE",
      "DistanceKm": 373227.05417
    }
  },
  "CurrentState": {
    "MoonDays": 23.03674,
    "Illumination": 43.78443,
    "Phase": {
      "Name": "Waning Crescent",
      "NameLocalized": "Убывающий серп",
      "Emoji": "🌘",
      "IsWaxing": false
    },
    "Zodiac": {
      "Name": "Gemini",
      "NameLocalized": "Близнецы",
      "Emoji": "♊"
    },
    "MoonPosition": {
      "Time": "2025-09-15T14:00:00+05:00",
      "AzimuthDegrees": 279.37345,
      "AltitudeDegrees": 29.04833,
      "Direction": "W",
      "DistanceKm": 374869.93889
    }
  },
  "EndDay": {
    "MoonDays": 23.53674,
    "Illumination": 38.2726,
    "Phase": {
      "Name": "Waning Crescent",
      "NameLocalized": "Убывающий серп",
      "Emoji": "🌘",
      "IsWaxing": false
    },
    "Zodiac": {
      "Name": "Gemini",
      "NameLocalized": "Близнецы",
      "Emoji": "♊"
    },
    "MoonPosition": {
      "Time": "2025-09-16T02:00:00+05:00",
      "AzimuthDegrees": 57.08873,
      "AltitudeDegrees": 8.79646,
      "Direction": "ENE",
      "DistanceKm": 376559.88437
    }
  },
  "MoonDaysDetailed": {
    "Count": 2,
    "Day": [
      {
        "Begin": "2025-09-14T21:31:05+05:00",
        "IsBeginExists": true,
        "End": "2025-09-15T22:37:45+05:00",
        "IsEndExists": true
      },
      {
        "Begin": "2025-09-15T22:37:45+05:00",
        "IsBeginExists": true,
        "End": "2025-09-16T23:56:10+05:00",
        "IsEndExists": true
      }
    ]
  },
  "ZodiacDetailed": {
    "Count": 1,
    "Zodiac": [
      {
        "Name": "Gemini",
        "NameLocalized": "Близнецы",
        "Emoji": "♊",
        "Begin": "2025-09-14T23:07:06+05:00",
        "End": "2025-09-17T11:07:06+05:00"
      }
    ]
  },
  "MoonRiseAndSet": {
    "IsMoonRise": true,
    "IsMoonSet": true,
    "IsMeridian": true,
    "Moonrise": {
      "Time": "2025-09-15T22:37:45+05:00",
      "AzimuthDegrees": 42.31555,
      "AltitudeDegrees": -0.56667,
      "Direction": "NE",
      "DistanceKm": 376365.00012
    },
    "Moonset": {
      "Time": "2025-09-15T15:46:53+05:00",
      "AzimuthDegrees": 318.44328,
      "AltitudeDegrees": -0.56667,
      "Direction": "NW",
      "DistanceKm": 375398.22022
    },
    "Meridian": {
      "Time": "2025-09-15T06:40:13+05:00",
      "AzimuthDegrees": 180,
      "AltitudeDegrees": 67.1,
      "Direction": "S",
      "DistanceKm": 374133.37617
    }
  }
}
```

</details>


<details>
  <summary><strong>timeFormat = 2006-01-02 15:04</strong></summary>

```json
{
    "BeginDay": {
        "MoonDays": 0.85568,
        "Illumination": 0.35744,
        "Phase": {
            "Name": "New Moon",
            "NameLocalized": "Новолуние",
            "Emoji": "🌑",
            "IsWaxing": true
        },
        "Zodiac": {
            "Name": "Virgo",
            "NameLocalized": "Дева",
            "Emoji": "♍"
        },
        "MoonPosition": {
            "Time": "2025-01-01 00:00",
            "AzimuthDegrees": 326.86255,
            "AltitudeDegrees": -62.68617,
            "Direction": "NNW",
            "DistanceKm": 382464.56722
        }
    },
    "CurrentState": {
        "MoonDays": 0.89806,
        "Illumination": 0.4137,
        "Phase": {
            "Name": "New Moon",
            "NameLocalized": "Новолуние",
            "Emoji": "🌑",
            "IsWaxing": true
        },
        "Zodiac": {
            "Name": "Virgo",
            "NameLocalized": "Дева",
            "Emoji": "♍"
        },
        "MoonPosition": {
            "Time": "2025-01-01 01:01",
            "AzimuthDegrees": 356.08557,
            "AltitudeDegrees": -65.49632,
            "Direction": "N",
            "DistanceKm": 382316.55173
        }
    },
    "EndDay": {
        "MoonDays": 1.85568,
        "Illumination": 2.79491,
        "Phase": {
            "Name": "New Moon",
            "NameLocalized": "Новолуние",
            "Emoji": "🌑",
            "IsWaxing": true
        },
        "Zodiac": {
            "Name": "Virgo",
            "NameLocalized": "Дева",
            "Emoji": "♍"
        },
        "MoonPosition": {
            "Time": "2025-01-02 00:00",
            "AzimuthDegrees": 309.21002,
            "AltitudeDegrees": -54.25535,
            "Direction": "NW",
            "DistanceKm": 379194.7976
        }
    },
    "MoonDaysDetailed": {
        "Count": 2,
        "Day": [
            {
                "Begin": "2024-12-31 09:18",
                "IsBeginExists": true,
                "End": "2025-01-01 09:55",
                "IsEndExists": true
            },
            {
                "Begin": "2025-01-01 09:55",
                "IsBeginExists": true,
                "End": "2025-01-02 10:21",
                "IsEndExists": true
            }
        ]
    },
    "ZodiacDetailed": {
        "Count": 1,
        "Zodiac": [
            {
                "Name": "Virgo",
                "NameLocalized": "Дева",
                "Emoji": "♍",
                "Begin": "2024-12-31 03:27",
                "End": "2025-01-02 15:27"
            }
        ]
    },
    "MoonRiseAndSet": {
        "IsMoonRise": true,
        "IsMoonSet": true,
        "IsMeridian": true,
        "Moonrise": {
            "Time": "2025-01-01 09:55",
            "AzimuthDegrees": 133.52155,
            "AltitudeDegrees": -0.56667,
            "Direction": "SE",
            "DistanceKm": 381054.53085
        },
        "Moonset": {
            "Time": "2025-01-01 17:26",
            "AzimuthDegrees": 228.55471,
            "AltitudeDegrees": -0.56667,
            "Direction": "SW",
            "DistanceKm": 380040.1844
        },
        "Meridian": {
            "Time": "2025-01-01 13:36",
            "AzimuthDegrees": 180,
            "AltitudeDegrees": 13.2,
            "Direction": "S",
            "DistanceKm": 380550.46947
        }
    }
}
```

</details>

<details>
  <summary><strong>timeFormat = timestamp</strong></summary>

```json
{
    "BeginDay": {
        "MoonDays": 0.85568,
        "Illumination": 0.35744,
        "Phase": {
            "Name": "New Moon",
            "NameLocalized": "Новолуние",
            "Emoji": "🌑",
            "IsWaxing": true
        },
        "Zodiac": {
            "Name": "Virgo",
            "NameLocalized": "Дева",
            "Emoji": "♍"
        },
        "MoonPosition": {
            "Time": 1735671600,
            "AzimuthDegrees": 326.86255,
            "AltitudeDegrees": -62.68617,
            "Direction": "NNW",
            "DistanceKm": 382464.56722
        }
    },
    "CurrentState": {
        "MoonDays": 0.89806,
        "Illumination": 0.4137,
        "Phase": {
            "Name": "New Moon",
            "NameLocalized": "Новолуние",
            "Emoji": "🌑",
            "IsWaxing": true
        },
        "Zodiac": {
            "Name": "Virgo",
            "NameLocalized": "Дева",
            "Emoji": "♍"
        },
        "MoonPosition": {
            "Time": 1735675261,
            "AzimuthDegrees": 356.08557,
            "AltitudeDegrees": -65.49632,
            "Direction": "N",
            "DistanceKm": 382316.55173
        }
    },
    "EndDay": {
        "MoonDays": 1.85568,
        "Illumination": 2.79491,
        "Phase": {
            "Name": "New Moon",
            "NameLocalized": "Новолуние",
            "Emoji": "🌑",
            "IsWaxing": true
        },
        "Zodiac": {
            "Name": "Virgo",
            "NameLocalized": "Дева",
            "Emoji": "♍"
        },
        "MoonPosition": {
            "Time": 1735758000,
            "AzimuthDegrees": 309.21002,
            "AltitudeDegrees": -54.25535,
            "Direction": "NW",
            "DistanceKm": 379194.7976
        }
    },
    "MoonDaysDetailed": {
        "Count": 2,
        "Day": [
            {
                "Begin": 1735618722,
                "IsBeginExists": true,
                "End": 1735707321,
                "IsEndExists": true
            },
            {
                "Begin": 1735707321,
                "IsBeginExists": true,
                "End": 1735795276,
                "IsEndExists": true
            }
        ]
    },
    "ZodiacDetailed": {
        "Count": 1,
        "Zodiac": [
            {
                "Name": "Virgo",
                "NameLocalized": "Дева",
                "Emoji": "♍",
                "Begin": 1735597669,
                "End": 1735813669
            }
        ]
    },
    "MoonRiseAndSet": {
        "IsMoonRise": true,
        "IsMoonSet": true,
        "IsMeridian": true,
        "Moonrise": {
            "Time": 1735707321,
            "AzimuthDegrees": 133.52155,
            "AltitudeDegrees": -0.56667,
            "Direction": "SE",
            "DistanceKm": 381054.53085
        },
        "Moonset": {
            "Time": 1735734376,
            "AzimuthDegrees": 228.55471,
            "AltitudeDegrees": -0.56667,
            "Direction": "SW",
            "DistanceKm": 380040.1844
        },
        "Meridian": {
            "Time": 1735720605,
            "AzimuthDegrees": 180,
            "AltitudeDegrees": 13.2,
            "Direction": "S",
            "DistanceKm": 380550.46947
        }
    }
}
```

</details>

----------------------------------------------------------------

### GET /api/v1/moonPhaseCurrent

The method returns the Moon parameters for the current day and time. If the day or time is not specified, the current value for the unspecified fields is taken. If longitude and latitude are specified, the response will contain additional structures.

This is a synonym for the moonPhaseDate method without day and time Params.

#### Params

  | Parameter | Type | Description | Example Value |
| :--- | :--- | :--- |  :--- | 
|`utc` | `string [optional, default="UTC+0"]` | UTC in format `UTC+7`, `UTC+09:30`, `-3` | `UTC+4`
|`lang` | `string [optional, default="en"]` | Values available: ("en", "es", "fr", "de", "ru", "jp") | `es`
|`precision` | `int [optional, default=2]` | How many digits after ```.``` will be in output. Allowed range: [1, 20] | `5`
|`latitude` | `float [optional, default=none]` | Latitude of viewer's place. Used for moon position calculations: ```MoonDaysDetailed```, ```MoonRiseAndSet```, and ```MoonPosition``` object | `51.1655`
|`longitude` | `float [optional, default=none]` | Longitude of viewer's place. Used for moon position calculations: ```MoonDaysDetailed```, ```MoonRiseAndSet```, and ```MoonPosition``` object | `71.4272`
|`timeFormat` | `string [optional]` | How the data will be formed. Available values: ISO (default), timestamp, duration or Go format. Duration is seconds between current date and event time.| `duration`

#### Response

Response as [GET /api/v1/moonPhaseDate](https://github.com/prostraction/moon/#v1moonphasedate-response)

----------------------------------------------------------------

### GET /api/v1/moonPhaseTimestamp

The method returns the Moon parameters for the given timestamp. If it is not specified, the current value for the timestamp is taken. If longitude and latitude are specified, the response will contain additional structures.

This is a synonym for the moonPhaseDate method but with timestamp instead of date.

#### Params

  | Parameter | Type | Description | Example Value |
| :--- | :--- | :--- |  :--- | 
|`utc` | `string [optional, default="UTC+0"]` | UTC in format `UTC+7`, `UTC+09:30`, `-3` | `UTC+4`
|`lang` | `string [optional, default="en"]` | Values available: ("en", "es", "fr", "de", "ru", "jp") | `es`
|`precision` | `int [optional, default=2]` | How many digits after ```.``` will be in output. Allowed range: [1, 20] | `5`
|`latitude` | `float [optional, default=none]` | Latitude of viewer's place. Used for moon position calculations: ```MoonDaysDetailed```, ```MoonRiseAndSet```, and ```MoonPosition``` object | `51.1655`
|`longitude` | `float [optional, default=none]` | Longitude of viewer's place. Used for moon position calculations: ```MoonDaysDetailed```, ```MoonRiseAndSet```, and ```MoonPosition``` object | `71.4272`
|`timestamp` | `int [optional, default=<current>]` | Timestamp for calculations | `1758045697`
|`timeFormat` | `string [optional]` | How the data will be formed. Available values: ISO (default), timestamp, duration or Go format. Duration is seconds between current date and event time.| `duration`

#### Response

Response as [GET /api/v1/moonPhaseDate](https://github.com/prostraction/moon/#v1moonphasedate-response)

----------------------------------------------------------------

### GET /api/v1/moonPositionMonthly

The method returns Moon position for specified month.

#### Params

  | Parameter | Type | Description | Example Value |
| :--- | :--- | :--- |  :--- | 
|`utc` | `string [optional, default="UTC+0"]` | UTC in format `UTC+7`, `UTC+09:30`, `-3` | `UTC+4`
|`precision` | `int [optional, default=2]` | How many digits after ```.``` will be in output. Allowed range: [1, 20] | `5`
|`latitude` | `float [required]` | Latitude of viewer's place. Used for moon position calculations: ```MoonDaysDetailed```, ```MoonRiseAndSet```, and ```MoonPosition``` object | `51.1655`
|`longitude` | `float [required]` | Longitude of viewer's place. Used for moon position calculations: ```MoonDaysDetailed```, ```MoonRiseAndSet```, and ```MoonPosition``` object | `71.4272`
|`year` | `int [optional, default=<current year>]` | Format: YYYY Allowed range: [1, 9999] | `2025`
|`month` | `int [optional, default=<current month>]` | Format: M or MM. Allowed range: [1, 12] | `01` or `1`
|`timeFormat` | `string [optional]` | How the data will be formed. Available values: ISO (default), timestamp, duration or Go format. Duration is seconds between current date and event time.| `duration`

#### Response

The method returns array of object MoonRiseAndSet for each day of selected month.

##### MoonRiseAndSet

```MoonRiseAndSet``` is a structure for determining the moonrise, moonset and meridian on a given day. Exists only if latitude and longitude are specified.

<details>
  <summary><strong>Table</strong></summary>

| Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`MoonRiseAndSet.Date` | `String [optional]` | Date for day of calculations, missing as default | `2025-01-15` |
|`MoonRiseAndSet.IsMoonRise` | `Boolean [required]` | True if moonrise occurs at given day | `true` |
|`MoonRiseAndSet.IsMoonSet` | `Boolean [required]` | True if moonset occurs at given day | `true` |
|`MoonRiseAndSet.IsMeridian` | `Boolean [required]` | True if meridian transit occurs at given day | `true` |
|`MoonRiseAndSet.Moonrise` | `Object of struct MoonPosition [optional]` | Moonrise position data. Exists only if IsMoonRise = true | - |
|`MoonRiseAndSet.Moonset` | `Object of struct MoonPosition [optional]` | Moonset position data. Exists only if IsMoonSet = true | - |
|`MoonRiseAndSet.Meridian` | `Object of struct MoonPosition [optional]` | Meridian position data, Exists only if IsMeridian = true | - |

MoonPosition structure:

<details>
  <summary><strong>Table</strong></summary>

| Response Variable | Type | Description | Example Value |
| :--- | :--- | :--- | :--- |
|`MoonPosition.Time` | `String / int64 [required]` | Time in ISO (default) or specified timeFormat | `"2025-09-16T00:00:00+05:00"` |
|`MoonPosition.AzimuthDegrees` | `Float [required]` | Moonrise azimuth | `47.3` |
|`MoonPosition.AltitudeDegrees` | `Float [required]` | Moonrise altitude | `-0.6` |
|`MoonPosition.Direction` | `String [required]` | Moonrise direction | `"ENE"` |
|`MoonPosition.DistanceKm` | `Float [required]` | Earth-Moon distance in km | `376559.9` |

</details>

</details>

#### Response example

<details>
  <summary><strong>JSON</strong></summary>

```json
[
    {
        "Date": "2023-03-01",
        "IsMoonRise": true,
        "IsMoonSet": true,
        "IsMeridian": true,
        "Moonrise": {
            "Time": "2023-03-01T11:07:37+05:00",
            "AzimuthDegrees": 42.8762,
            "AltitudeDegrees": -0.56667,
            "Direction": "NE",
            "DistanceKm": 402636.75966
        },
        "Moonset": {
            "Time": "2023-03-01T04:10:07+05:00",
            "AzimuthDegrees": 316.74613,
            "AltitudeDegrees": -0.56667,
            "Direction": "NW",
            "DistanceKm": 401798.89266
        },
        "Meridian": {
            "Time": "2023-03-01T20:05:30+05:00",
            "AzimuthDegrees": 180,
            "AltitudeDegrees": 66.2,
            "Direction": "S",
            "DistanceKm": 403572.05151
        }
    },
    {
        "Date": "2023-03-02",
        "IsMoonRise": true,
        "IsMoonSet": true,
        "IsMeridian": true,
        "Moonrise": {
            "Time": "2023-03-02T12:01:07+05:00",
            "AzimuthDegrees": 43.38104,
            "AltitudeDegrees": -0.56667,
            "Direction": "NE",
            "DistanceKm": 404836.13007
        },
        "Moonset": {
            "Time": "2023-03-02T05:03:00+05:00",
            "AzimuthDegrees": 317.04922,
            "AltitudeDegrees": -0.56667,
            "Direction": "NW",
            "DistanceKm": 404345.126
        },
        "Meridian": {
            "Time": "2023-03-02T20:56:36+05:00",
            "AzimuthDegrees": 180,
            "AltitudeDegrees": 65.4,
            "Direction": "S",
            "DistanceKm": 405326.66664
        }
    },
    ...
    {
        "Date": "2023-03-31",
        "IsMoonRise": true,
        "IsMoonSet": true,
        "IsMeridian": true,
        "Moonrise": {
            "Time": "2023-03-31T11:58:33+05:00",
            "AzimuthDegrees": 49.55137,
            "AltitudeDegrees": -0.56667,
            "Direction": "NE",
            "DistanceKm": 404935.96867
        },
        "Moonset": {
            "Time": "2023-03-31T04:17:06+05:00",
            "AzimuthDegrees": 312.1322,
            "AltitudeDegrees": -0.56667,
            "Direction": "NW",
            "DistanceKm": 404804.48186
        },
        "Meridian": {
            "Time": "2023-03-31T20:28:24+05:00",
            "AzimuthDegrees": 180,
            "AltitudeDegrees": 61.6,
            "Direction": "S",
            "DistanceKm": 404936.05398
        }
    }
]
```

</details>

<details>
  <summary><strong>timeFormat '2006-01-02 15:04'</strong></summary>

```json
[
    {
        "Date": "2023-03-01",
        "IsMoonRise": true,
        "IsMoonSet": true,
        "IsMeridian": true,
        "Moonrise": {
            "Time": "2023-03-01 11:07",
            "AzimuthDegrees": 42.8762,
            "AltitudeDegrees": -0.56667,
            "Direction": "NE",
            "DistanceKm": 402636.75966
        },
        "Moonset": {
            "Time": "2023-03-01 04:10",
            "AzimuthDegrees": 316.74613,
            "AltitudeDegrees": -0.56667,
            "Direction": "NW",
            "DistanceKm": 401798.89266
        },
        "Meridian": {
            "Time": "2023-03-01 20:05",
            "AzimuthDegrees": 180,
            "AltitudeDegrees": 66.2,
            "Direction": "S",
            "DistanceKm": 403572.05151
        }
    },
 ...
```

</details>


<details>
  <summary><strong>timeFormat 'timestamp'</strong></summary>

```json
[
    {
        "Date": "2023-03-01",
        "IsMoonRise": true,
        "IsMoonSet": true,
        "IsMeridian": true,
        "Moonrise": {
            "Time": 1677650857,
            "AzimuthDegrees": 42.8762,
            "AltitudeDegrees": -0.56667,
            "Direction": "NE",
            "DistanceKm": 402636.75966
        },
        "Moonset": {
            "Time": 1677625807,
            "AzimuthDegrees": 316.74613,
            "AltitudeDegrees": -0.56667,
            "Direction": "NW",
            "DistanceKm": 401798.89266
        },
        "Meridian": {
            "Time": 1677683130,
            "AzimuthDegrees": 180,
            "AltitudeDegrees": 66.2,
            "Direction": "S",
            "DistanceKm": 403572.05151
        }
    },
 ...
```
</details>

----------------------------------------------------------------

### GET /api/v1/moonTableYear

The method returns the moon phases for the given year. The response contains an array for each month, each element of which contains the time of the new moon, first quarter, full moon, last quarter.

#### Params

  | Parameter | Type | Description | Example Value |
| :--- | :--- | :--- |  :--- | 
|`utc` | `string [optional, default="UTC+0"]` | UTC in format `UTC+7`, `UTC+09:30`, `-3` | `UTC+4`
|`year` | `int [optional, default=<current year>]` | Format: YYYY Allowed range: [1, 9999] | `2025`
|`timeFormat` | `string [optional]` | How the data will be formed. Available values: ISO (default), timestamp, duration or Go format. Duration is seconds between current date and event time.| `duration`

#### Response example

<details>
  <summary><strong>JSON</strong></summary>

```json
[
    // first moon of the year
     {
        "NewMoon": "2024-12-31T08:27:49+10:00",
        "FirstQuarter": "2025-01-07T16:15:21+10:00",
        "FullMoon": "2025-01-14T08:27:44+10:00",
        "LastQuarter": "2025-01-22T13:25:30+10:00"
    },
    // second moon of the year
    {
        "NewMoon": "2025-01-29T22:37:18+10:00",
        "FirstQuarter": "2025-02-06T09:58:10+10:00",
        "FullMoon": "2025-02-12T23:54:26+10:00",
        "LastQuarter": "2025-02-21T20:11:33+10:00"
    },
...
    // last moon of the year
     {
        "NewMoon": "2025-12-20T11:44:25+10:00",
        "FirstQuarter": "2025-12-27T19:53:55+10:00",
        "FullMoon": "2026-01-03T20:04:15+10:00",
        "LastQuarter": "2026-01-11T02:44:20+10:00"
    }
]
```

</details>

<details>
  <summary><strong>timeFormat '2006-01-02 15:04'</strong></summary>

```json
[
    // first moon of the year
    {
        "NewMoon": "2022-12-23 15:17",
        "FirstQuarter": "2022-12-30 06:24",
        "FullMoon": "2023-01-07 04:09",
        "LastQuarter": "2023-01-07 04:09"
    },
    // second moon of the year
    {
        "NewMoon": "2023-01-22 01:55",
        "FirstQuarter": "2023-01-28 20:05",
        "FullMoon": "2023-02-05 23:30",
        "LastQuarter": "2023-02-05 23:30"
    },
...
    // last moon of the year
    {
        "NewMoon": "2023-12-13 04:32",
        "FirstQuarter": "2023-12-19 23:42",
        "FullMoon": "2023-12-27 05:33",
        "LastQuarter": "2023-12-27 05:33"
    }
]
```

</details>

<details>
  <summary><strong>timeFormat 'timestamp'</strong></summary>

```json
[
    // first moon of the year
    {
        "NewMoon": 1671790676,
        "FirstQuarter": 1672363444,
        "FullMoon": 1673046595,
        "LastQuarter": 1673046595
    },
    // second moon of the year
    {
        "NewMoon": 1674334530,
        "FirstQuarter": 1674918315,
        "FullMoon": 1675621844,
        "LastQuarter": 1675621844
    },
...
    // last moon of the year
    {
        "NewMoon": 1702423927,
        "FirstQuarter": 1703011321,
        "FullMoon": 1703637223,
        "LastQuarter": 1703637223
    }
]
```

</details>

----------------------------------------------------------------

### GET /api/v1/moonTableCurrent

The method returns the moon phases for the current year. The response contains an array for each month, each element of which contains the time of the new moon, first quarter, full moon, last quarter.

#### Params

  | Parameter | Type | Description | Example Value |
| :--- | :--- | :--- |  :--- | 
|`utc` | `string [optional, default="UTC+0"]` | UTC in format `UTC+7`, `UTC+09:30`, `-3` | `UTC+4`
|`timeFormat` | `string [optional]` | How the data will be formed. Available values: ISO (default), timestamp, duration or Go format. Duration is seconds between current date and event time.| `duration`

#### Response

Response: as [GET /api/v1/moonTableYear](https://github.com/prostraction/moon#v1moontableyear-response)

----------------------------------------------------------------

### GET /api/v1/toJulianTimeByDate

The method converts human date to julian time (UTC +0 timezone).

#### Params

  | Parameter | Type | Description | Example Value |
| :--- | :--- | :--- |  :--- | 
|`precision` | `int [optional, default=2]` | How many digits after ```.``` will be in output. Allowed range: [1, 20] | `5`
|`timeFormat` | `string [optional]` | How the data will be formed. Available values: ISO (default), timestamp, duration or Go format. Duration is seconds between current date and event time.| `duration`
|`year` | `int [optional, default=<current year>]` | Format: YYYY Allowed range: [1, 9999] | `2025`
|`month` | `int [optional, default=<current month>]` | Format: M or MM. Allowed range: [1, 12] | `01` or `1`
|`day` | `int [optional, default=<current day>]` | Format: D or DD. Allowed range: [1, 31] | `01` or `1`
|`hour` | `int [optional, default=<current hour>]` | Format: h or hh. Allowed range: [0, 23] | `01` or `1`
|`minute` | `int [optional, default=<current minute>]` | Format: m or mm. Allowed range: [0, 59] | `01` or `1`
|`second` | `int [optional, default=<current second>]` | Format: s or ss. Allowed range: [0, 59] | `01` or `1`

#### Response examples

```json
{
    "CivilDate": "2025-01-01T01:01:01Z",
    "JulianDate": 2460676.54237
}
```

timeFormat: `Mon, 02 Jan 2006 15:04:05`

```json
{
    "CivilDate": "Wed, 01 Jan 2025 01:01:01",
    "JulianDate": 2460676.54237
}
```

timeFormat: `timestamp`

```json
{
    "CivilDate": 1735693261,
    "JulianDate": 2460676.54237
}
```

----------------------------------------------------------------

### GET /api/v1/toJulianTimeByTimestamp

The method converts human date to julian time (UTC +0 timezone).

#### Params

  | Parameter | Type | Description | Example Value |
| :--- | :--- | :--- |  :--- | 
|`precision` | `int [optional, default=2]` | How many digits after ```.``` will be in output. Allowed range: [1, 20] | `5`
|`timestamp` | `int [optional, default=<current>]` | Timestamp for calculations. Current, if not specified | `1735693261`
|`timeFormat` | `string [optional]` | How the data will be formed. Available values: ISO (default), timestamp, duration or Go format. Duration is seconds between current date and event time.| `duration`

#### Response example

```json
{
    "CivilDate": "2025-01-01T01:01:01Z",
    "JulianDate": 2460676.54237
}
```

timeFormat: `Mon, 02 Jan 2006 15:04:05`

```json
{
    "CivilDate": "Wed, 01 Jan 2025 01:01:01",
    "JulianDate": 2460676.54237
}
```

timeFormat: `timestamp`

```json
{
    "CivilDate": 1735693261,
    "JulianDate": 2460676.54237
}

----------------------------------------------------------------

### GET /api/v1/fromJulianTime

The method converts julian time to human time (UTC +0 timezone).

#### Params

  | Parameter | Type | Description | Example Value |
| :--- | :--- | :--- |  :--- | 
|`precision` | `int [optional, default=2]` | How many digits after ```.``` will be in output. Allowed range: [1, 20] | `5`
|`jtime` | `float64 [required]` | Julian Time to convert (float64) | `2460676.5423726854`
|`timeFormat` | `string [optional]` | How the data will be formed. Available values: ISO (default), timestamp, duration or Go format. Duration is seconds between current date and event time.| `duration`

#### Response examples

```json
{
    "CivilDate": "2025-01-01T01:01:01Z",
    "JulianDate": 2460676.54237
}
```

timeFormat: `Mon, 02 Jan 2006 15:04:05`

```json
{
    "CivilDate": "Wed, 01 Jan 2025 01:01:01",
    "JulianDate": 2460676.54237
}
```

timeFormat: `timestamp`

```json
{
    "CivilDate": 1735693261,
    "JulianDate": 2460676.54237
}
```

### GET /api/v1/nextMoonPhase

The method returns next moon phase

#### Params

  | Parameter | Type | Description | Example Value |
| :--- | :--- | :--- |  :--- | 
|`utc` | `string [optional, default="UTC+0"]` | UTC in format `UTC+7`, `UTC+09:30`, `-3` | `UTC+4`
|`timeFormat` | `string [optional]` | How the data will be formed. Available values: ISO (default), timestamp, duration or Go format. Duration is seconds between current date and event time.| `duration`


#### Response example:

timeFormat: `ISO`

```json
{
    "NewMoon": "2025-11-20T11:48:16+05:00",
    "FirstQuarter": "2025-11-28T11:44:00+05:00",
    "FullMoon": "2025-12-05T04:15:18+05:00",
    "LastQuarter": "2025-12-12T02:04:31+05:00"
}
```

timeFormat: `Mon, 02 Jan 2006 15:04:05`

```json
{
    "NewMoon": "Thu, 20 Nov 2025 11:48:16",
    "FirstQuarter": "Fri, 28 Nov 2025 11:44:00",
    "FullMoon": "Fri, 05 Dec 2025 04:15:18",
    "LastQuarter": "Fri, 12 Dec 2025 02:04:31"
}
```

timeFormat: `timestamp`

```json
{
    "NewMoon": 1763621296,
    "FirstQuarter": 1764312240,
    "FullMoon": 1764890118,
    "LastQuarter": 1765487071
}
```

timeFormat: `duration`

```json
{
    "NewMoon": 313830,
    "FirstQuarter": 1004774,
    "FullMoon": 1582652,
    "LastQuarter": 2179605
}
```
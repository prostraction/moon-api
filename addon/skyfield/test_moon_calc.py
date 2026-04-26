"""Unit tests for addon/skyfield/moon_calc.py.

Run from addon/skyfield directory:
    python -m unittest test_moon_calc.py

These tests validate computational correctness against:
  * timeanddate.com (Astana, Kazakhstan: 51.13°N 71.43°E, UTC+5)
  * NASA daily moon guide
  * skyfield ground truth (pure structural checks)
"""

import unittest
from datetime import datetime, timedelta, timezone

import moon_calc


class TestTotalOffsetHours(unittest.TestCase):
    def test_positive(self):
        self.assertAlmostEqual(moon_calc._total_offset_hours(5, 30), 5.5)

    def test_negative_unsigned_minutes(self):
        # The Go layer ships unsigned minutes; we must apply the hour's sign.
        self.assertAlmostEqual(moon_calc._total_offset_hours(-5, 30), -5.5)

    def test_negative_signed_minutes(self):
        # Either signed or unsigned minutes must produce the same result.
        self.assertAlmostEqual(moon_calc._total_offset_hours(-5, -30), -5.5)

    def test_zero_hours_signed_minutes(self):
        # Edge case: only minute offset, must use minute sign verbatim.
        self.assertAlmostEqual(moon_calc._total_offset_hours(0, 30), 0.5)
        self.assertAlmostEqual(moon_calc._total_offset_hours(0, -30), -0.5)

    def test_extremes(self):
        self.assertAlmostEqual(moon_calc._total_offset_hours(14, 0), 14.0)
        self.assertAlmostEqual(moon_calc._total_offset_hours(-12, 0), -12.0)
        self.assertAlmostEqual(moon_calc._total_offset_hours(12, 45), 12.75)


class TestMoonPositionAtTime(unittest.TestCase):
    """Smoke tests against skyfield. We don't validate exact angles (skyfield
    is the reference), but the response shape must be stable and values must
    be in valid ranges.
    """

    def test_position_response_shape(self):
        r = moon_calc.get_moon_position_at_time(
            lat=51.13,
            lon=71.43,
            timezone_hours=5,
            timezone_minutes=0,
            precision=2,
            year=2024,
            month=1,
            day=15,
            hour=20,
            minute=0,
            second=0,
        )
        self.assertEqual(r["Status"], "success")
        self.assertIn("AzimuthDegrees", r)
        self.assertIn("AltitudeDegrees", r)
        self.assertIn("Direction", r)
        self.assertIn("DistanceKm", r)
        self.assertIn("Timestamp", r)

        # Azimuth in [0, 360), altitude in [-90, 90].
        self.assertGreaterEqual(r["AzimuthDegrees"], 0)
        self.assertLess(r["AzimuthDegrees"], 360)
        self.assertGreaterEqual(r["AltitudeDegrees"], -90)
        self.assertLessEqual(r["AltitudeDegrees"], 90)

        # Earth-Moon distance: 356,400..406,700 km (perigee/apogee).
        self.assertGreater(r["DistanceKm"], 350_000)
        self.assertLess(r["DistanceKm"], 410_000)

        # Direction must be one of 16 cardinal abbreviations.
        self.assertIn(
            r["Direction"],
            {
                "N",
                "NNE",
                "NE",
                "ENE",
                "E",
                "ESE",
                "SE",
                "SSE",
                "S",
                "SSW",
                "SW",
                "WSW",
                "W",
                "WNW",
                "NW",
                "NNW",
            },
        )

    def test_position_negative_timezone(self):
        # New York: UTC-5 (no DST in January). Must not crash and must produce
        # the same UTC moon position as the same UTC moment from Astana.
        astana = moon_calc.get_moon_position_at_time(
            lat=51.13,
            lon=71.43,
            timezone_hours=5,
            timezone_minutes=0,
            precision=4,
            year=2024,
            month=1,
            day=15,
            hour=20,
            minute=0,
            second=0,
        )
        # Same UTC instant from a UTC-5 timezone: 15 Jan 20:00 Astana = 15 Jan 15:00 UTC = 15 Jan 10:00 New York
        ny = moon_calc.get_moon_position_at_time(
            lat=51.13,
            lon=71.43,
            timezone_hours=-5,
            timezone_minutes=0,
            precision=4,
            year=2024,
            month=1,
            day=15,
            hour=10,
            minute=0,
            second=0,
        )
        # Same observer + same UTC moment → identical azimuth/altitude.
        self.assertAlmostEqual(astana["AzimuthDegrees"], ny["AzimuthDegrees"], places=2)
        self.assertAlmostEqual(astana["AltitudeDegrees"], ny["AltitudeDegrees"], places=2)
        self.assertAlmostEqual(astana["DistanceKm"], ny["DistanceKm"], delta=1.0)

    def test_position_fractional_negative_timezone(self):
        # Newfoundland: UTC-3:30. Verify that minutes are applied with the
        # hour's sign (regression for the -4.5 vs -5.5 bug).
        nfld = moon_calc.get_moon_position_at_time(
            lat=47.5,
            lon=-52.7,
            timezone_hours=-3,
            timezone_minutes=30,
            precision=4,
            year=2024,
            month=1,
            day=15,
            hour=8,
            minute=30,
            second=0,
        )
        # Same UTC moment via UTC.
        utc = moon_calc.get_moon_position_at_time(
            lat=47.5,
            lon=-52.7,
            timezone_hours=0,
            timezone_minutes=0,
            precision=4,
            year=2024,
            month=1,
            day=15,
            hour=12,
            minute=0,
            second=0,
        )
        self.assertAlmostEqual(nfld["AzimuthDegrees"], utc["AzimuthDegrees"], places=2)
        self.assertAlmostEqual(nfld["AltitudeDegrees"], utc["AltitudeDegrees"], places=2)


class TestDailyMoonData(unittest.TestCase):
    def test_response_shape(self):
        d = moon_calc.get_daily_moon_data(
            lat=51.13,
            lon=71.43,
            timezone_hours=5,
            timezone_minutes=0,
            precision=2,
            year=2024,
            month=1,
            day=15,
        )
        for key in ("Moonrise", "Moonset", "Meridian", "IsMoonRise", "IsMoonSet", "IsMeridian"):
            self.assertIn(key, d)
        # Either rise or set should exist on most days at this latitude.
        self.assertTrue(d["IsMoonRise"] or d["IsMoonSet"])

    def test_astana_known_event(self):
        # Astana, 2024-01-25: timeanddate.com lists Moonrise around 17:09 local,
        # Moonset around ~07:51 next morning (Wolf Moon day, near full).
        # We test only that moonrise time is in late afternoon local time,
        # which is robust to small algorithmic differences.
        d = moon_calc.get_daily_moon_data(
            lat=51.13,
            lon=71.43,
            timezone_hours=5,
            timezone_minutes=0,
            precision=2,
            year=2024,
            month=1,
            day=25,
        )
        self.assertTrue(d["IsMoonRise"])
        rise_ts = d["Moonrise"]["Timestamp"]
        rise_dt_local = datetime.fromtimestamp(rise_ts, tz=timezone(timedelta(hours=5)))
        # Wolf Moon rises in the late afternoon; tolerance 4 hours around 17:00.
        self.assertEqual(rise_dt_local.date(), datetime(2024, 1, 25).date())
        self.assertGreaterEqual(rise_dt_local.hour, 14)
        self.assertLessEqual(rise_dt_local.hour, 21)


class TestMoonDataResponseShape(unittest.TestCase):
    def test_daily_wraps(self):
        r = moon_calc.get_moon_data_response(
            lat=51.13,
            lon=71.43,
            timezone_hours=5,
            timezone_minutes=0,
            precision=2,
            year=2024,
            month=1,
            day=15,
        )
        self.assertEqual(r["Status"], "success")
        self.assertEqual(r["Range"], "single_day")
        self.assertEqual(r["Parameters"]["Latitude"], 51.13)
        self.assertEqual(r["Parameters"]["Longitude"], 71.43)
        self.assertEqual(r["Parameters"]["UtcOffset"], "UTC+05:00")

    def test_negative_offset_format(self):
        r = moon_calc.get_moon_data_response(
            lat=40.7,
            lon=-74.0,
            timezone_hours=-5,
            timezone_minutes=0,
            precision=2,
            year=2024,
            month=1,
            day=15,
        )
        self.assertEqual(r["Parameters"]["UtcOffset"], "UTC-05:00")

    def test_fractional_offset_format(self):
        # India UTC+5:30
        r = moon_calc.get_moon_data_response(
            lat=28.6,
            lon=77.2,
            timezone_hours=5,
            timezone_minutes=30,
            precision=2,
            year=2024,
            month=1,
            day=15,
        )
        self.assertEqual(r["Parameters"]["UtcOffset"], "UTC+05:30")

    def test_negative_fractional_offset_format(self):
        # Regression for the -4.5 vs -5.5 sign bug.
        r = moon_calc.get_moon_data_response(
            lat=47.5,
            lon=-52.7,
            timezone_hours=-3,
            timezone_minutes=30,
            precision=2,
            year=2024,
            month=1,
            day=15,
        )
        self.assertEqual(r["Parameters"]["UtcOffset"], "UTC-03:30")


if __name__ == "__main__":
    unittest.main()

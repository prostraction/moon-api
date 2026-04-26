"""Integration tests for addon/skyfield/server.py.

Spawns the HTTP server on a random port and exercises HTTP error paths.

Run from addon/skyfield directory:
    python -m unittest test_server.py
"""

import json
import threading
import time
import unittest
import urllib.request
from http.server import HTTPServer

import server as srv


def _free_port():
    import socket

    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.bind(("127.0.0.1", 0))
        return s.getsockname()[1]


class _ServerThread:
    def __init__(self):
        self.port = _free_port()
        self.httpd = HTTPServer(("127.0.0.1", self.port), srv.MoonRequestHandler)
        # silence stderr by reusing the handler's silent log_message
        self.thread = threading.Thread(target=self.httpd.serve_forever, daemon=True)

    def __enter__(self):
        self.thread.start()
        # tiny pause to make sure socket is listening
        time.sleep(0.05)
        return self

    def __exit__(self, exc_type, exc, tb):
        self.httpd.shutdown()
        self.httpd.server_close()
        self.thread.join(timeout=2)

    def get(self, path):
        url = f"http://127.0.0.1:{self.port}{path}"
        try:
            with urllib.request.urlopen(url, timeout=10) as r:
                return r.status, json.loads(r.read().decode("utf-8"))
        except urllib.error.HTTPError as e:
            body = e.read().decode("utf-8")
            try:
                return e.code, json.loads(body)
            except json.JSONDecodeError:
                return e.code, body


class TestServerErrorHandling(unittest.TestCase):
    def test_unknown_path_404(self):
        with _ServerThread() as s:
            status, body = s.get("/foo")
            self.assertEqual(status, 404)
            self.assertEqual(body["Status"], "error")
            self.assertIn("unknown path", body["Message"])

    def test_invalid_lat_400(self):
        with _ServerThread() as s:
            status, body = s.get("/position?lat=abc&lon=0")
            self.assertEqual(status, 400)
            self.assertEqual(body["Status"], "error")

    def test_lat_out_of_range_400(self):
        with _ServerThread() as s:
            status, body = s.get("/position?lat=95&lon=0")
            self.assertEqual(status, 400)
            self.assertIn("lat out of range", body["Message"])

    def test_lon_out_of_range_400(self):
        with _ServerThread() as s:
            status, body = s.get("/position?lat=0&lon=200")
            self.assertEqual(status, 400)
            self.assertIn("lon out of range", body["Message"])

    def test_utc_hours_out_of_range_400(self):
        with _ServerThread() as s:
            status, body = s.get("/position?lat=51&lon=71&utc_hours=24")
            self.assertEqual(status, 400)
            self.assertIn("utc_hours out of range", body["Message"])

    def test_utc_minutes_invalid(self):
        with _ServerThread() as s:
            status, body = s.get("/position?lat=51&lon=71&utc_minutes=70")
            self.assertEqual(status, 400)

    def test_invalid_month(self):
        with _ServerThread() as s:
            status, body = s.get("/daily?lat=51&lon=71&utc_hours=5&year=2024&month=13&day=1")
            self.assertEqual(status, 400)
            self.assertIn("month", body["Message"])

    def test_position_success(self):
        with _ServerThread() as s:
            status, body = s.get("/position?lat=51.13&lon=71.43&utc_hours=5&year=2024&month=1&day=25&hour=20")
            self.assertEqual(status, 200)
            self.assertEqual(body["Status"], "success")
            self.assertIn("AzimuthDegrees", body)

    def test_daily_success(self):
        with _ServerThread() as s:
            status, body = s.get("/daily?lat=51.13&lon=71.43&utc_hours=5&year=2024&month=1&day=25")
            self.assertEqual(status, 200)
            self.assertEqual(body["Status"], "success")
            self.assertIn("Data", body)

    def test_options_preflight(self):
        with _ServerThread() as s:
            req = urllib.request.Request(f"http://127.0.0.1:{s.port}/position", method="OPTIONS")
            with urllib.request.urlopen(req, timeout=5) as r:
                self.assertEqual(r.status, 200)
                self.assertIn("Access-Control-Allow-Origin", dict(r.headers))


if __name__ == "__main__":
    unittest.main()

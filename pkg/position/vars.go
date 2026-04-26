package position

import "time"

// baseURL is a var (not const) so tests can substitute an httptest.Server.
var baseURL = "http://localhost:9997/"

// httpClientTimeout must stay below the server WriteTimeout (60s) so we surface
// upstream failures before the parent request times out.
const httpClientTimeout = 30 * time.Second

// SetBaseURLForTesting overrides the upstream base URL and returns a function
// that restores the previous value. Intended for tests in other packages that
// need to mock the upstream service.
func SetBaseURLForTesting(url string) func() {
	prev := baseURL
	baseURL = url
	return func() { baseURL = prev }
}

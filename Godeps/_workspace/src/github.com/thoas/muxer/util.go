package muxer

import (
	"net/http"
	"strings"
)

// matchMap returns true if the given key/value pairs exist in a given map.
func matchMap(toCheck map[string]string, toMatch map[string][]string,
	canonicalKey bool) bool {
	for k, v := range toCheck {
		// Check if key exists.
		if canonicalKey {
			k = http.CanonicalHeaderKey(k)
		}
		if values := toMatch[k]; values == nil {
			return false
		} else if v != "" {
			// If value was defined as an empty string we only check that the
			// key exists. Otherwise we also check for equality.
			valueExists := false
			for _, value := range values {
				if v == value {
					valueExists = true
					break
				}
			}
			if !valueExists {
				return false
			}
		}
	}
	return true
}

func getHost(r *http.Request) string {
	if r.URL.IsAbs() {
		return r.URL.Host
	}
	host := r.Host
	// Slice off any port information.
	if i := strings.Index(host, ":"); i != -1 {
		host = host[:i]
	}
	return host

}

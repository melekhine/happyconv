package webconv

import (
	"log"
	"net/http"
)

// Log colors
const (
	endString    = " %0.2f kB \u2192 %s"
	InfoColor    = "\033[1;34m%s\033[0m" + endString
	WarningColor = "\033[1;33m%s\033[0m" + endString
)

// WebconvLogMiddleware logs http requests
func WebconvLogMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		color := InfoColor

		if r.Method != "POST" {
			color = WarningColor
		}

		log.Printf(color, r.Method, float64(r.ContentLength)/1024.0, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

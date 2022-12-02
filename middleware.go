package barter

import (
	"context"
	"mime"
	"net/http"
	"strings"

	"google.golang.org/api/idtoken"
)

// ValidateAcceptHeader validates the request has a supported "Accept" header
// MIME type. If the Accept header is invalid a 406 is returned and the
// request is terminated.
func ValidateAcceptHeader(validAcceptTypes map[string]struct{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acceptVal := r.Header.Get("Accept")
			for _, val := range strings.Split(acceptVal, ",") {
				t, _, err := mime.ParseMediaType(val)
				if err != nil {
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}

				_, exists := validAcceptTypes[t]
				if !exists {
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ValidateJWT validates the JWT in the "Authorization" header. If the JWT
// is missing or invalid a 401 is returned and the request is terminated.
// If the JWT is valid the trimmed token is stored in the "Token-raw"
// header and the subject is stored in the "Subject" header.
func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		tokenElements := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(tokenElements) < 2 {
			RenderError(w, r, http.StatusUnauthorized, err, "")
			return
		}
		payload, err := idtoken.Validate(context.TODO(), tokenElements[1], "")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			RenderError(w, r, http.StatusUnauthorized, err, "")
			return
		}
		r.Header.Set("Token-raw", tokenElements[1])
		r.Header.Set("Subject", payload.Subject)
		next.ServeHTTP(w, r)
	})
}

// MustExtractSubject extracts the subject from the JWT in the Authorization header.
// If validation fails (the JWT is invalid or missing) the subject header is set
// to an empty string.
func MustExtractSubject(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		tokenElements := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(tokenElements) < 2 {
			r.Header.Set("Subject", "")
			next.ServeHTTP(w, r)
			return
		}
		payload, err := idtoken.Validate(context.TODO(), tokenElements[1], "")
		if err != nil {
			r.Header.Set("Subject", "")
		} else {
			r.Header.Set("Subject", payload.Subject)
		}
		next.ServeHTTP(w, r)
	})
}

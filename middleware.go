package barter

import (
	"context"
	"mime"
	"net/http"
	"strings"

	"google.golang.org/api/idtoken"
)

// ValidateAcceptHeader validates the request has a supported Accept header MIME type.
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

// Validate extracts and validates the subject from the JWT in the Authorization
// header. If validation fails (the JWT is invalid or missing) the request
// cascade is terminated and a 401 is returned. If validation succeeds the
// "Subject" header of the request is set to the subject extracted from the JWT
// and the "Token-raw" header of the request is set to the raw token.
func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		tokenElements := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(tokenElements) < 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("missing jwt"))
			return
		}
		payload, err := idtoken.Validate(context.TODO(), tokenElements[1], "")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("error validating: " + err.Error()))
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

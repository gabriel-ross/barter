package barter

import (
	"context"
	"net/http"
	"strings"

	"google.golang.org/api/idtoken"
)

// Validate extracts and validates the subject from the JWT in the Authorization
// header. If validation fails (the JWT is invalid or missing) the request
// cascade is terminated and a 401 is returned. If validation succeeds the
// "Subject" header of the request is set to the subject extracted from the JWT.
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
		r.Header.Set("Token", tokenElements[1])
		r.Header.Set("Subject", payload.Subject)
		next.ServeHTTP(w, r)
	})
}

// ExtractSubject extracts the subject from the JWT in the Authorization header.
// If validation fails (the JWT is invalid or missing) the subject header is set
// to an empty string.
func ExtractSubject(next http.Handler) http.Handler {
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

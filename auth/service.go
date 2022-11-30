package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
)

type Service struct {
	router       chi.Router
	db           *firestore.Client
	absolutePath string
	codeSlug     string
	tokenSlug    string
	oauth2Config *oauth2.Config
}

// New mounts the auth service router at baseURL+"/"+resourceSlug and returns
// a pointer to the new service. By default code and token endpoints are mounted
// to /code and /token, respectively. Both of these URLs must be added to the
// OAuth2.0 provider's list of permitted redirect URLs.
func New(router chi.Router, db *firestore.Client, oauth2Config *oauth2.Config, baseURL, serviceSlug string, options ...Option) *Service {
	svc := &Service{
		router:       router,
		db:           db,
		absolutePath: fmt.Sprintf("%s/%s", baseURL, serviceSlug),
		codeSlug:     "code",
		tokenSlug:    "token",
		oauth2Config: oauth2Config,
	}

	for _, opt := range options {
		opt(svc)
	}

	router.Mount("/"+serviceSlug, svc.Routes())
	return svc
}

type Option func(*Service)

func WithCodeSlug(slug string) func(*Service) {
	return func(svc *Service) {
		svc.codeSlug = slug
	}
}

func WithTokenSlug(slug string) func(*Service) {
	return func(svc *Service) {
		svc.tokenSlug = slug
	}
}

func RandString(length int) (_ string, err error) {
	t := make([]byte, length)
	_, err = rand.Read(t)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(t), nil
}

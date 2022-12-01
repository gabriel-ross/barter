package user

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
)

type Service struct {
	router                 chi.Router
	db                     *firestore.Client
	absolutePath           string
	supportedResponseMIMEs map[string]struct{}
}

// New mounts the User service router at baseURL+"/"+resourceSlug and returns
// a pointer to the new service.
func New(router chi.Router, db *firestore.Client, baseURL, resourceSlug string, supportedResponseMIMEs map[string]struct{}) *Service {
	svc := &Service{
		router:                 router,
		db:                     db,
		absolutePath:           fmt.Sprintf("%s/%s", baseURL, resourceSlug),
		supportedResponseMIMEs: supportedResponseMIMEs,
	}
	router.Mount("/"+resourceSlug, svc.Routes())
	return svc
}

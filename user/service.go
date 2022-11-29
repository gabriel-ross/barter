package user

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
)

type Service struct {
	router     chi.Router
	db         *firestore.Client
	absResPath string
}

// New mounts the User service router at baseURL+"/"+resourceSlug and returns
// a point to the new service.
func New(router chi.Router, db *firestore.Client, baseURL, resourceSlug string) *Service {
	svc := &Service{
		router:     router,
		db:         db,
		absResPath: fmt.Sprintf("%s/%s", baseURL, resourceSlug),
	}
	router.Mount("/"+resourceSlug, svc.Routes())
	return svc
}

package transaction

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
)

type Service struct {
	router       chi.Router
	db           *firestore.Client
	absolutePath string
}

// New mounts the Transaction service router at baseURL+"/"+resourceSlug and returns
// a point to the new service.
func New(router chi.Router, db *firestore.Client, baseURL, resourceSlug string) *Service {
	svc := &Service{
		router:       router,
		db:           db,
		absolutePath: fmt.Sprintf("%s/%s", baseURL, resourceSlug),
	}
	router.Mount("/"+resourceSlug, svc.Routes())
	return svc
}

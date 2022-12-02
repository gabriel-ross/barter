package user

import (
	"github.com/gabriel-ross/barter"
	"github.com/go-chi/chi"
)

func (svc *Service) Routes() chi.Router {
	r := chi.NewRouter()

	r.With(barter.ValidateAcceptHeader(svc.supportedResponseMIMEs)).Get("/", svc.handleList())

	r.Route("/{id}", func(r chi.Router) {
		r.Use(barter.ValidateJWT, svc.ValidateUserExistsAndRequestorAccess)

		r.With(barter.ValidateAcceptHeader(svc.supportedResponseMIMEs)).Get("/", svc.handleGet())
		r.Put("/", svc.handlePut())
		r.Patch("/", svc.handlePatch())
		r.Delete("/", svc.handleDelete())
	})

	return r
}

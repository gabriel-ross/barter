package account

import (
	"github.com/gabriel-ross/barter"
	"github.com/go-chi/chi"
)

func (svc *Service) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(barter.ValidateJWT)

	r.With(barter.ValidateAcceptHeader(svc.supportedResponseMIMEs)).Post("/", svc.handleCreate())
	r.With(barter.ValidateAcceptHeader(svc.supportedResponseMIMEs)).Get("/", svc.handleList())

	r.Route("/{id}", func(r chi.Router) {
		r.Use(svc.ValidateAccountExistsAndRequestorAccess)

		r.With(barter.ValidateAcceptHeader(svc.supportedResponseMIMEs)).Get("/", svc.handleGet())
		r.Put("/", svc.handlePut())
		r.Patch("/", svc.handlePatch())
		r.Delete("/", svc.handleDelete())

		r.Route("/owner", func(r chi.Router) {
			r.Put("/", svc.setOwner())
			r.Delete("/", svc.removeOwner())
		})
	})

	return r
}

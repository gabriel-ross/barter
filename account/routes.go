package account

import (
	"github.com/gabriel-ross/barter"
	"github.com/go-chi/chi"
)

func (svc *Service) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(barter.ValidateJWT, svc.ValidateAccountExistsAndRequestorAccess)

	r.With(barter.ValidateAcceptHeader(svc.supportedResponseMIMEs)).Post("/", svc.handleCreate())
	r.With(barter.ValidateAcceptHeader(svc.supportedResponseMIMEs)).Get("/", svc.handleList())
	r.Route("/{id}", func(r chi.Router) {
		r.With(barter.ValidateAcceptHeader(svc.supportedResponseMIMEs)).Get("/", svc.handleGet())
		r.Put("/", svc.handleUpdate())
		r.Delete("/", svc.handleDelete())

		r.Route("/user", func(r chi.Router) {
			r.Put("/", svc.setUser())
			r.Delete("/", svc.removeUser())
		})
	})

	return r
}

package user

import (
	"github.com/gabriel-ross/barter"
	"github.com/go-chi/chi"
)

func (svc *Service) Routes() chi.Router {
	r := chi.NewRouter()

	r.With(barter.ValidateAcceptHeader(svc.supportedResponseMIMEs)).Post("/", svc.handleCreate())
	r.With(barter.ValidateAcceptHeader(svc.supportedResponseMIMEs)).Get("/", svc.handleList())
	r.Route("/{id}", func(r chi.Router) {
		r.With(barter.ValidateAcceptHeader(svc.supportedResponseMIMEs), svc.ValidateUserExistsAndRequestorAccess).Get("/", svc.handleGet())
		r.With(svc.ValidateUserExistsAndRequestorAccess).Put("/", svc.handleUpdate())
		r.With(svc.ValidateUserExistsAndRequestorAccess).Delete("/", svc.handleDelete())
	})

	return r
}

package transaction

import (
	"github.com/gabriel-ross/barter"
	"github.com/go-chi/chi"
)

func (svc *Service) Routes() chi.Router {
	r := chi.NewRouter()

	r.With(barter.ValidateAcceptHeader(svc.supportedResponseMIMEs)).Post("/", svc.handleCreate())
	r.With(barter.ValidateAcceptHeader(svc.supportedResponseMIMEs)).Get("/", svc.handleList())

	r.Route("/{id}", func(r chi.Router) {
		r.Use(svc.VerifyTransactionExists)

		r.With(barter.ValidateAcceptHeader(svc.supportedResponseMIMEs)).Get("/", svc.handleGet())
		r.Put("/", svc.handlePut())
		r.Patch("/", svc.handlePatch())
		r.Delete("/", svc.handleDelete())

		r.Route("/sender", func(r chi.Router) {
			r.Put("/", svc.setSender())
			r.Delete("/", svc.removeSender())
		})

		r.Route("/recipient", func(r chi.Router) {
			r.Put("/", svc.setRecipient())
			r.Delete("/", svc.removeRecipient())
		})
	})

	return r
}

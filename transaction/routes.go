package transaction

import "github.com/go-chi/chi"

func (svc *Service) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", svc.handleCreate())
	r.Get("/", svc.handleList())
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", svc.handleGet())
		r.Put("/", svc.handleUpdate())
		r.Delete("/", svc.handleDelete())
	})

	return r
}

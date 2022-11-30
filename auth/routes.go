package auth

import "github.com/go-chi/chi"

func (svc *Service) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/"+svc.codeSlug, svc.getCode())
	r.Get("/"+svc.tokenSlug, svc.getToken())

	return r
}

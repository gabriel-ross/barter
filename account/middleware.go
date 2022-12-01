package account

import (
	"context"
	"errors"
	"net/http"

	"github.com/gabriel-ross/barter"
	"github.com/go-chi/chi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VerifyAccountExists verifies the account with id "id" exists. If it does
// not it terminates the request and returns 404.
func (svc *Service) VerifyAccountExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error

		id := chi.URLParam(r, "id")
		_, err = svc.read(ctx, id)
		if status.Code(err) == codes.NotFound {
			barter.RenderError(w, r, http.StatusNotFound, errors.New("resource not found"), "%s", err.Error())
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ValidateAccountExistsAndRequestorAccess validates the account with "id" exists
// and that the requestor has permission to access the resource. If the account has
// no owner any user id can modify it.
func (svc *Service) ValidateAccountExistsAndRequestorAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error

		id := chi.URLParam(r, "id")
		data, err := svc.read(ctx, id)
		if status.Code(err) == codes.NotFound {
			barter.RenderError(w, r, http.StatusNotFound, errors.New("resource not found"), "%s", err.Error())
			return
		}

		resourceOwner := data.UserID
		requestor := r.Header.Get("Subject")
		if resourceOwner != "" && requestor != resourceOwner {
			barter.RenderError(w, r, http.StatusForbidden, nil, "user %s does not have permission to modify resource owned by %s", requestor, resourceOwner)
			return
		}

		next.ServeHTTP(w, r)
	})
}

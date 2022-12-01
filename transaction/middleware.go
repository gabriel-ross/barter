package transaction

import (
	"context"
	"errors"
	"net/http"

	"github.com/gabriel-ross/barter"
	"github.com/go-chi/chi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VerifyTransactionExists verifies the Transaction with id "id" exists. If it does
// not it terminates the request and returns 404.
func (svc *Service) VerifyTransactionExists(next http.Handler) http.Handler {
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

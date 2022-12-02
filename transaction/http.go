package transaction

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/barter"
	"github.com/gabriel-ross/barter/model"
	"github.com/go-chi/chi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (svc *Service) handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		data := model.NewTransaction()

		err = BindRequest(r, &data)
		if err != nil {
			barter.RenderError(w, r, http.StatusBadRequest, err, "%s", err.Error())
			return
		}

		resp, err := svc.create(ctx, data)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		svc.RenderResponse(w, r, http.StatusCreated, resp)
	}
}

func (svc *Service) handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error

		offset, limit, err := extractPaginate(r)
		if err != nil {
			barter.RenderError(w, r, http.StatusBadRequest, err, "%s", err.Error())
			return
		}

		resp, err := svc.list(ctx, barter.WithOrder("id", firestore.Asc), barter.WithOffset(offset), barter.WithLimit(limit))
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		count, err := svc.count(ctx)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		svc.RenderListResponse(w, r, http.StatusOK, resp, offset, limit, count)
	}
}

func (svc *Service) handleGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")

		resp, err := svc.read(ctx, id)
		if status.Code(err) == codes.NotFound {
			barter.RenderError(w, r, http.StatusNotFound, errors.New("resource not found"), "%s", err.Error())
			return
		} else if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		svc.RenderResponse(w, r, http.StatusOK, resp)
	}
}

func (svc *Service) handlePut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")
		data := model.NewTransaction()

		err = BindRequest(r, &data)
		if err != nil {
			barter.RenderError(w, r, http.StatusBadRequest, err, "%s", err.Error())
			return
		}

		_, err = svc.update(ctx, id, data)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (svc *Service) handlePatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("transaction patch"))
	}
}

func (svc *Service) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")

		_, err = svc.read(ctx, id)
		if status.Code(err) == codes.NotFound {
			barter.RenderError(w, r, http.StatusNotFound, errors.New("resource not found"), "%s", err.Error())
			return
		} else if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		err = svc.delete(ctx, id)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (svc *Service) setSender() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")

		m := model.NewTransaction()
		err = BindRequest(r, &m)

		data, err := svc.read(ctx, id)
		if err != nil && status.Code(err) != codes.NotFound {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		data.SenderAccountID = m.SenderAccountID
		_, err = svc.update(ctx, id, data)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (svc *Service) removeSender() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")

		data, err := svc.read(ctx, id)
		if err != nil && status.Code(err) != codes.NotFound {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		data.SenderAccountID = ""
		_, err = svc.update(ctx, id, data)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (svc *Service) setRecipient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")

		m := model.NewTransaction()
		err = BindRequest(r, &m)

		data, err := svc.read(ctx, id)
		if err != nil && status.Code(err) != codes.NotFound {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		data.RecipientAccountID = m.RecipientAccountID
		_, err = svc.update(ctx, id, data)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (svc *Service) removeRecipient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")

		data, err := svc.read(ctx, id)
		if err != nil && status.Code(err) != codes.NotFound {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		data.RecipientAccountID = ""
		_, err = svc.update(ctx, id, data)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func extractPaginate(r *http.Request) (_ int, _ int, err error) {
	offset := 0
	limit := 5

	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		offset, err = strconv.Atoi(offsetParam)
		if err != nil {
			return 0, 0, err
		}
	}

	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return 0, 0, err
		}
	}

	return offset, limit, nil
}

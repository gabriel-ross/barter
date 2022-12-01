package account

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
		req := model.NewAccount()

		err = BindRequest(r, &req)
		if err != nil {
			barter.RenderError(w, r, http.StatusBadRequest, err, "%s", err.Error())
			return
		}

		data := model.NewAccount()
		data.UserID = req.UserID
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

		requestor := r.Header.Get("Subject")
		resp, err := svc.list(ctx, barter.WithFilter("userID", barter.Eq, requestor), barter.WithOrder("id", firestore.Asc), barter.WithOffset(offset), barter.WithLimit(limit))
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		count, err := svc.count(ctx, barter.WithFilter("userID", barter.Eq, requestor))
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

func (svc *Service) handleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")
		data := model.NewAccount()

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

func (svc *Service) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")

		_, err = svc.read(ctx, id)
		if err != nil && status.Code(err) != codes.NotFound {
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

func (svc *Service) setUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")

		req := model.NewAccount()
		err = BindRequest(r, &req)

		data, err := svc.read(ctx, id)
		if err != nil && status.Code(err) != codes.NotFound {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		data.UserID = req.UserID
		_, err = svc.update(ctx, id, data)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (svc *Service) removeUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")

		data, err := svc.read(ctx, id)
		if err != nil && status.Code(err) != codes.NotFound {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		data.UserID = ""
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

package user

import (
	"context"
	"errors"
	"net/http"
	"strconv"

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
		data := model.User{}

		err = BindRequest(r, &data)
		if err != nil {
			barter.RenderError(w, r, http.StatusBadRequest, err)
			return
		}

		resp, err := svc.create(ctx, data)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err)
			return
		}

		svc.RenderResponse(w, r, resp, http.StatusCreated)
	}
}

func (svc *Service) handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error

		offset, limit, err := extractPaginate(r)
		if err != nil {
			barter.RenderError(w, r, http.StatusBadRequest, err)
			return
		}
		resp, err := svc.list(ctx, offset, limit)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err)
			return
		}

		count, err := svc.count(ctx)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err)
			return
		}

		svc.RenderListResponse(w, r, http.StatusOK, resp, offset, limit, count)
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

func (svc *Service) handleGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")

		resp, err := svc.read(ctx, id)
		if status.Code(err) == codes.NotFound {
			barter.RenderError(w, r, http.StatusNotFound, errors.New("resource not found"))
			return
		} else if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err)
			return
		}

		svc.RenderResponse(w, r, resp, http.StatusOK)
	}
}

func (svc *Service) handleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")
		data := model.User{}

		err = BindRequest(r, &data)
		if err != nil {
			w.Write([]byte("error binding: " + err.Error()))
			return
		}

		resp, err := svc.update(ctx, id, data)
		if err != nil {
			w.Write([]byte("error creating: " + err.Error()))
			return
		}

		svc.RenderResponse(w, r, resp, http.StatusNoContent)
	}
}

func (svc *Service) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")

		exists, err := svc.exists(ctx, id)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err)
			return
		}

		if !exists {
			barter.RenderError(w, r, http.StatusNotFound, errors.New("resource not found"))
			return
		}

		err = svc.delete(ctx, id)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}
}

package barter

import (
	"encoding/json"
	"net/http"
)

// Renderer is an interface
type Renderer interface {
	Render(*http.Request) any
}

func Render(w http.ResponseWriter, r *http.Request, code int, data Renderer) {
	var err error
	respData := data.Render(r)
	respBody, err := json.Marshal(respData)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(code)
	w.Write(respBody)
	return
}
func RenderList(w http.ResponseWriter, r *http.Request, code int, data []Renderer) {}

func RenderError(w http.ResponseWriter, r *http.Request, code int, svrErr error) {
	var err error
	errResp := newErrorResponse(code, svrErr)
	respBody, err := json.Marshal(errResp)
	if err != nil {
		writeErrorError(w, err)
		return
	}

	w.WriteHeader(code)
	w.Write(respBody)
	return
}

type errorResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"-"`
	AppCode        int64  `json:"code,omitempty"`
	ErrorText      string `json:"error,omitempty"`
}

func newErrorResponse(code int, err error) *errorResponse {
	return &errorResponse{
		Err:            err,
		HTTPStatusCode: code,
		ErrorText:      err.Error(),
	}
}

func writeErrorError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("error encountered while attempting to write error: " + err.Error()))
	return
}

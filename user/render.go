package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gabriel-ross/barter"
	"github.com/gabriel-ross/barter/model"
)

type response struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Self        string `json:"self"`
}

type listResponse struct {
	Data   []response `json:"data"`
	Count  int        `json:"count"`
	Offset int        `json:"offset"`
	Limit  int        `json:"limit"`
	Prev   string     `json:"prev,omitempty"`
	Next   string     `json:"next,omitempty"`
	Self   string     `json:"self"`
}

func (svc *Service) newResponse(u model.User) response {
	return response{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Self:        fmt.Sprintf("%s/%s", svc.absResPath, u.ID),
	}
}

func (svc *Service) RenderResponse(w http.ResponseWriter, r *http.Request, u model.User, code int) {
	resp := svc.newResponse(u)
	body, err := json.Marshal(resp)
	if err != nil {
		barter.RenderError(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(code)
	w.Write(body)
	return
}

func (svc *Service) RenderListResponse(w http.ResponseWriter, r *http.Request, code int, users []model.User, offset, limit, count int) {
	data := []response{}
	for _, user := range users {
		data = append(data, svc.newResponse(user))
	}

	resp := listResponse{
		Data:   data,
		Count:  count,
		Offset: offset,
		Limit:  limit,
		Self:   svc.absResPath,
	}

	if offset > 0 {
		resp.Prev = fmt.Sprintf("%s?offset=%d&limit=%d", svc.absResPath, max(0, offset-limit), limit)
	}
	if offset+limit < count {
		resp.Next = fmt.Sprintf("%s?offset=%d&limit=%d", svc.absResPath, offset+limit, limit)
	}

	body, err := json.Marshal(resp)
	if err != nil {
		barter.RenderError(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(code)
	w.Write(body)
	return
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

package transaction

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gabriel-ross/barter"
	"github.com/gabriel-ross/barter/model"
)

type response struct {
	ID                 string         `json:"id"`
	Quantities         map[string]int `json:"quantities"`
	SenderAccountID    string         `json:"sender"`
	RecipientAccountID string         `json:"recipient"`
	Timestamp          time.Time      `json:"timestamp"`
	Self               string         `json:"self"`
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

func (svc *Service) newResponse(m model.Transaction) response {
	return response{
		ID:                 m.ID,
		Timestamp:          m.Timestamp,
		Quantities:         m.Quantities,
		SenderAccountID:    m.SenderAccountID,
		RecipientAccountID: m.RecipientAccountID,
		Self:               fmt.Sprintf("%s/%s", svc.absolutePath, m.ID),
	}
}

func (svc *Service) RenderResponse(w http.ResponseWriter, r *http.Request, code int, m model.Transaction) {
	barter.WriteResponse(w, r, code, svc.newResponse(m))
}

func (svc *Service) RenderListResponse(w http.ResponseWriter, r *http.Request, code int, transactions []model.Transaction, offset, limit, count int) {
	data := []response{}
	for _, transaction := range transactions {
		data = append(data, svc.newResponse(transaction))
	}

	resp := listResponse{
		Data:   data,
		Count:  count,
		Offset: offset,
		Limit:  limit,
		Self:   svc.absolutePath,
	}

	if offset > 0 {
		resp.Prev = fmt.Sprintf("%s?offset=%d&limit=%d", svc.absolutePath, max(0, offset-limit), limit)
	}
	if offset+limit < count {
		resp.Next = fmt.Sprintf("%s?offset=%d&limit=%d", svc.absolutePath, offset+limit, limit)
	}

	barter.WriteResponse(w, r, code, resp)
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

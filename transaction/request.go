package transaction

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gabriel-ross/barter/model"
)

type request struct {
	Quantities         map[string]int `json:"quantities"`
	SenderAccountID    string         `json:"sender"`
	RecipientAccountID string         `json:"recipient"`
}

// BindRequest binds the fields defined in request of a request to a User.
// This method also extracts the token from the header "Token".
func BindRequest(r *http.Request, m *model.Transaction) (err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var reqBody request
	err = json.Unmarshal(body, &reqBody)

	m.Quantities = reqBody.Quantities
	m.SenderAccountID = reqBody.SenderAccountID
	m.RecipientAccountID = reqBody.RecipientAccountID
	m.Timestamp = time.Now()
	return nil
}

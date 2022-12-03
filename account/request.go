package account

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gabriel-ross/barter/model"
)

type request struct {
	Owner      string             `json:"owner"`
	Balances   map[string]float64 `json:"balances"`
	Reputation int                `json:"reputation"`
}

// BindRequest binds the fields defined in body of a request to an Account.
// This method also extracts the token from the header "Token".
func BindRequest(r *http.Request, m *model.Account) (err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var reqBody request
	err = json.Unmarshal(body, &reqBody)

	m.Owner = reqBody.Owner
	m.Balances = reqBody.Balances
	m.Reputation = reqBody.Reputation
	return nil
}

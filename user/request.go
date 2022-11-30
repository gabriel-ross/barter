package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gabriel-ross/barter/model"
)

type request struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

// BindRequest binds the fields defined in request of a request to a User.
// This method also extracts the token from the header "Token".
func BindRequest(r *http.Request, m *model.User) (err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var reqBody request
	err = json.Unmarshal(body, &reqBody)

	m.Name = reqBody.Name
	m.Email = reqBody.Email
	m.PhoneNumber = reqBody.PhoneNumber
	m.Token = r.Header.Get("Token")
	return nil
}

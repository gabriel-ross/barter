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
func BindRequest(r *http.Request, u *model.User) (err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var reqBody request
	err = json.Unmarshal(body, &reqBody)

	u.Name = reqBody.Name
	u.Email = reqBody.Email
	u.PhoneNumber = reqBody.PhoneNumber
	u.Token = r.Header.Get("Token")
	return nil
}

package model

type User struct {
	ID          string `firestore:"id"`
	Name        string `firestore:"name"`
	Email       string `firestore:"email"`
	PhoneNumber string `firestore:"phoneNumber"`
	Token       string `firestore:"token"`
}

func NewUser() User {
	return User{
		ID:          "",
		Name:        "",
		Email:       "",
		PhoneNumber: "",
		Token:       "",
	}
}

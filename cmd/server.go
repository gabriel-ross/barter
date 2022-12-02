package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/barter"
	"github.com/gabriel-ross/barter/account"
	"github.com/gabriel-ross/barter/auth"
	"github.com/gabriel-ross/barter/transaction"
	"github.com/gabriel-ross/barter/user"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var APPLICATION_URL string
var PROJECT_ID string
var PORT string
var CLIENT_ID string
var CLIENT_SECRET string
var TOKEN_REDIRECT_URL string

func main() {
	var err error
	LoadConfigFromEnvironment()
	ctx := context.TODO()

	// Instantiate dependencies
	r := chi.NewRouter()
	fsClient, err := firestore.NewClient(ctx, PROJECT_ID)
	if err != nil {
		log.Fatalf("Failed to create client: %s", err)
	}
	defer fsClient.Close()

	oauth2Config := &oauth2.Config{
		ClientID:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"openid"},
		RedirectURL:  fmt.Sprintf("%s/auth/token", APPLICATION_URL),
	}

	// Mount index page
	r.Get("/", index())

	supportedMIMETypes := map[string]struct{}{
		"*/*":              struct{}{},
		"application/json": struct{}{},
	}

	// Instantiate and mount services
	auth.New(r, fsClient, oauth2Config, APPLICATION_URL, "auth")
	user.New(r, fsClient, APPLICATION_URL, "users", supportedMIMETypes)
	transaction.New(r, fsClient, APPLICATION_URL, "transactions", supportedMIMETypes)
	account.New(r, fsClient, APPLICATION_URL, "accounts", supportedMIMETypes)

	log.Println("Server up and running on port: ", PORT)
	http.ListenAndServe(PORT, r)
}

func LoadConfigFromEnvironment() {
	godotenv.Load(".env")
	APPLICATION_URL = os.Getenv("APPLICATION_URL")
	PROJECT_ID = os.Getenv("PROJECT_ID")
	PORT = ":" + os.Getenv("PORT")
	CLIENT_ID = os.Getenv("CLIENT_ID")
	CLIENT_SECRET = os.Getenv("CLIENT_SECRET")
	TOKEN_REDIRECT_URL = os.Getenv("TOKEN_REDIRECT_URL")

	// Default value if not set
	if PORT == "" {
		PORT = ":8080"
	}
}

type Config struct {
	PROJECT_ID string `env:"PROJECT_ID" required:"true" default:"-"`
	PORT       string `env:"PORT" required:"false" default:"8080"`
}

func index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		barter.WriteResponse(w, r, http.StatusOK, "hello world")
	}
}

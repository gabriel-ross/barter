package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/barter"
	"github.com/go-chi/chi"
	"google.golang.org/api/iterator"
)

var PROJECT_ID string
var PORT string

type demo struct {
	id string
}

func main() {
	ctx := context.TODO()
	LoadConfigFromEnvironment()
	db, err := barter.NewFirestoreClient(ctx, PROJECT_ID)
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Get("/", index)
	r.Post("/tests/{foo}", create(ctx, db, "tests"))
	r.Get("/tests", list(ctx, db, "tests"))
	r.Get("/tests/nofilter", listNoFilter(ctx, db, "tests"))
	r.Get("/tests/filteroffsetlimit", listLimitOffset(ctx, db, "tests"))
	r.Get("/tests/filter/{foo}", listFilter(ctx, db, "tests"))
	r.Get("/tests/filterall/{foo}", listFilterAll(ctx, db, "tests"))
	http.ListenAndServe(PORT, r)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("POST /tests create\nGET /tests list"))
}

func create(ctx context.Context, db *firestore.Client, collectionPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := db.Collection(collectionPath).NewDoc().ID
		data := map[string]interface{}{
			"id":  id,
			"foo": chi.URLParam(r, "foo"),
		}
		_, err := db.Collection(collectionPath).Doc(id).Set(ctx, data)
		if err != nil {
			w.Write([]byte("error storing data: " + err.Error()))
			return
		}

		bytes, err := json.Marshal(data)
		if err != nil {
			w.Write([]byte("error marshaling response: " + err.Error()))
			return
		}
		w.Write(bytes)
		return
	}
}

func list(ctx context.Context, db *firestore.Client, collectionPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		iter := db.Collection(collectionPath).Documents(ctx)
		resp := []map[string]interface{}{}
		for {
			dsnap, err := iter.Next()
			if err == iterator.Done {
				break
			}

			var data map[string]interface{}
			dsnap.DataTo(&data)
			resp = append(resp, data)
		}
		bytes, err := json.Marshal(resp)
		if err != nil {
			w.Write([]byte("error marshaling list: " + err.Error()))
			return
		}
		w.Write(bytes)
		return
	}
}

func listNoFilter(ctx context.Context, db *firestore.Client, collectionPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := barter.List[map[string]interface{}](ctx, db, collectionPath)
		if err != nil {
			w.Write([]byte("error retrieving data: " + err.Error()))
			return
		}
		bytes, err := json.Marshal(resp)
		if err != nil {
			w.Write([]byte("error marshaling list: " + err.Error()))
			return
		}
		w.Write(bytes)
		return
	}
}

func listFilter(ctx context.Context, db *firestore.Client, collectionPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := barter.List[map[string]interface{}](ctx, db, collectionPath, barter.WithFilterQuery("foo", barter.Eq, chi.URLParam(r, "foo")))
		if err != nil {
			w.Write([]byte("error retrieving data: " + err.Error()))
			return
		}
		bytes, err := json.Marshal(resp)
		if err != nil {
			w.Write([]byte("error marshaling list: " + err.Error()))
			return
		}
		w.Write(bytes)
		return
	}
}

func listFilterAll(ctx context.Context, db *firestore.Client, collectionPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := barter.List[map[string]interface{}](ctx, db, collectionPath, barter.WithFilterQuery("foo", barter.Eq, chi.URLParam(r, "foo")), barter.WithLimit(2), barter.WithOffset(1))
		if err != nil {
			w.Write([]byte("error retrieving data: " + err.Error()))
			return
		}
		bytes, err := json.Marshal(resp)
		if err != nil {
			w.Write([]byte("error marshaling list: " + err.Error()))
			return
		}
		w.Write(bytes)
		return
	}
}

func listLimitOffset(ctx context.Context, db *firestore.Client, collectionPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := barter.List[map[string]interface{}](ctx, db, collectionPath, barter.WithLimit(2), barter.WithOffset(1))
		if err != nil {
			w.Write([]byte("error retrieving data: " + err.Error()))
			return
		}
		bytes, err := json.Marshal(resp)
		if err != nil {
			w.Write([]byte("error marshaling list: " + err.Error()))
			return
		}
		w.Write(bytes)
		return
	}
}

func LoadConfigFromEnvironment() {
	//godotenv.Load(".env")
	PROJECT_ID = os.Getenv("PROJECT_ID")
	PORT = os.Getenv("PORT")

	// Default value if not set
	if PORT == "" {
		PORT = "8080"
	}

	PORT = ":" + PORT
}

type Config struct {
	PROJECT_ID string `env:"PROJECT_ID" required:"true" default:"-"`
	PORT       string `env:"PORT" required:"false" default:"8080"`
}

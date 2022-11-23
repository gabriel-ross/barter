package barter

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
)

type Storer interface {
	Exists(context.Context, string, string) (bool, error)
	Set(context.Context, string, string, interface{}) (interface{}, error)
}

type Issuer struct {
	router chi.Router
	path   string
	config *oauth2.Config
	db     Storer
}

func NewIssuer(path string, db Storer, clientID, clientSecret, redirectURL string, endpoints oauth2.Endpoint) *Issuer {
	newIss := &Issuer{
		router: chi.NewRouter(),
		db:     db,
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     endpoints,
			Scopes:       []string{"openid"},
			RedirectURL:  redirectURL,
		},
	}

	return newIss
}

func (i *Issuer) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/code", i.getCode())
	r.Get("/token", i.exchangeToken())
	return r
}

func (i *Issuer) getCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.TODO()
		state, err := i.genToken(30)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("unable to generate state: " + err.Error()))
			return
		}
		_, err = i.db.Set(ctx, "states", state, struct{}{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error storing state: " + err.Error()))
			return
		}
		http.Redirect(w, r, i.config.AuthCodeURL(state), http.StatusSeeOther)
	}
}

func (i *Issuer) exchangeToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.TODO()

		// Validate state
		stateIsValid, err := i.db.Exists(ctx, "states", r.FormValue("state"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error validating state: " + err.Error()))
			return
		}
		if !stateIsValid {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("invalid state: " + err.Error()))
			return
		}

		token, err := i.config.Exchange(ctx, r.FormValue("code"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error exchanging code %s: %s", r.FormValue("code"), err)))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token.Extra("id_token").(string)))
		return
	}
}

func (i *Issuer) genToken(length int) (_ string, err error) {
	t := make([]byte, length)
	_, err = rand.Read(t)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(t), nil
}

package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/gabriel-ross/barter"
	"github.com/gabriel-ross/barter/model"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (svc *Service) getCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.TODO()
		state, err := RandString(30)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "error generating state: %s", err.Error())
			return
		}
		err = svc.createState(ctx, state)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "error storing state: %s", err.Error())
			return
		}
		http.Redirect(w, r, svc.oauth2Config.AuthCodeURL(state), http.StatusSeeOther)
	}
}

func (svc *Service) getToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.TODO()

		// Validate state
		err = svc.read(ctx, r.FormValue("state"))
		if status.Code(err) == codes.NotFound {
			barter.RenderError(w, r, http.StatusBadRequest, err, "invalid state provided: %s\nerror: %s", r.FormValue("state"), err.Error())
		}
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "error encountered while validating state: %s", err.Error())
			return
		}

		token, err := svc.oauth2Config.Exchange(ctx, r.FormValue("code"))
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "error exchanging code %s\nerror: %s", r.FormValue("code"), err.Error())
			return
		}

		jwt := token.Extra("id_token").(string)
		payload, err := idtoken.Validate(ctx, jwt, "")
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "error encountered while validating token %s: %s", jwt, err.Error())
			return
		}

		log.Println("made it this far")

		// Register new user in database if they don't exist
		user := model.NewUser()
		user.ID = payload.Subject
		_, err = svc.createUser(ctx, user)
		if err != nil {
			barter.RenderError(w, r, http.StatusInternalServerError, err, "error encountered while creating new user %s: %s", jwt, err.Error())
			return
		}

		barter.WriteResponse(w, r, http.StatusOK, response{
			UserID: payload.Subject,
			Token:  jwt,
		})
	}
}

type response struct {
	UserID string `json:"userID"`
	Token  string `json:"token"`
}

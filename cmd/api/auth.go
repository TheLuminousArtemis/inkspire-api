package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/theluminousartemis/socialnews/internal/mailer"
	"github.com/theluminousartemis/socialnews/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=255"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Register a new user
//	@Description	This endpoint allows users to register by providing their username and password.
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	//hashed password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	//store the user

	plainToken := uuid.New().String()
	//store token
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err := app.storage.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestError(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	//mail
	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)
	isProdEnv := app.config.env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}
	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, isProdEnv)
	if err != nil {
		app.l.Errorw("error sending welcome error", "error", err)
		if err := app.storage.Users.Delete(ctx, user.ID); err != nil {
			app.l.Errorw("error deleting user", "error", err)
		}
		app.internalServerError(w, r, err)
		return
	}

	app.l.Infow("Email sent", "status code", status)

	if err = app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
	}
}

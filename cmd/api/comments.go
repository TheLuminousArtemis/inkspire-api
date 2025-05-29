package main

import (
	"net/http"

	"github.com/theluminousartemis/socialnews/internal/store"
)

type CommentPayload struct {
	Content string `json:"content" validate:"required"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CommentPayload

	post := getPostFromCtx(r)

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	comment := &store.Comment{
		Content: payload.Content,
		//change after auth
		PostID: post.ID,
		UserID: 2,
	}

	ctx := r.Context()
	if err := app.storage.Comments.Create(ctx, comment); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

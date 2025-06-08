package main

import (
	"net/http"

	"github.com/theluminousartemis/socialnews/internal/store"
)

type CommentPayload struct {
	Content string `json:"content" validate:"required"`
}

// CreateComment godoc
//
//	@Summary		Create a comment
//	@Description	Create a comment by ID
//	@Tags			posts, comments
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"Post ID"
//	@Param			payload	body		CommentPayload	true	"Comment payload"
//	@Success		200		{object}	store.Comment	"Comment created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id}/comments [post]
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
		PostID:  post.ID,
		//change after auth
		UserID: 2,
	}

	ctx := r.Context()
	if err := app.storage.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

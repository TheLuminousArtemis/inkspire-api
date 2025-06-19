package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/theluminousartemis/socialnews/internal/store"
)

type commentkey string

var commentCtxKey commentkey = "comment"

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
//	@Param			id		path		int								true	"Post ID"
//	@Param			payload	body		CommentPayload					true	"Comment payload"
//	@Success		200		{object}	store.SwaggerCommentResponse	"Comment created"
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

	user := getUserFromCtx(r)
	comment := &store.Comment{
		Content: payload.Content,
		PostID:  post.ID,
		UserID:  user.ID,
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

// DeleteComment godoc
//
//	@Summary		Delete a comment
//	@Description	Delete a comment by its ID if the user is the owner or has appropriate role
//	@Tags			posts, comments
//	@Accept			json
//	@Produce		json
//	@Param			postID		path		int		true	"Post ID"
//	@Param			commentID	path		int		true	"Comment ID"
//	@Success		204			{string}	string	"Comment deleted successfully"
//	@Failure		403			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID}/comments/{commentID} [delete]
func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	comment := getCommentfromCtx(r)
	ctx := r.Context()
	err := app.storage.Comments.Delete(ctx, comment.ID)
	if err != nil {
		app.internalServerError(w, r, err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) commentsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "commentID")
		commentID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		ctx := r.Context()
		comment, err := app.storage.Comments.GetByID(ctx, commentID)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.commentNotFoundErrorResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}
		ctx = context.WithValue(ctx, commentCtxKey, comment)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getCommentfromCtx(r *http.Request) *store.Comment {
	comment, _ := r.Context().Value(commentCtxKey).(*store.Comment)
	return comment
}

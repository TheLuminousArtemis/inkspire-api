package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.l.Errorf("internal server error: %v path: %s err: %v", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusInternalServerError, "Internal server error")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.l.Warnf("bad request error: %v path: %s err: %v", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.l.Errorf("bad request error: %v path: %s err: %v", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) userNotFoundErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.l.Warnf("user not found error: %v path: %s err: %v", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusNotFound, "user not found error")
}

func (app *application) postNotFoundErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.l.Warnf("post not found error: %v path: %s err: %v", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusNotFound, "Post not found error")
}

func (app *application) unauthorizedBasicResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.l.Warnf("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset=UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, "unauthorized request error")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.l.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.l.Warnw("forbidden response error", "method", r.Method, "path", r.URL.Path)
	writeJSONError(w, http.StatusForbidden, "forbidden")
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.l.Infow("rate limit exceeded", "remote_addr", r.RemoteAddr, "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Retry-After", retryAfter)
	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded")
}

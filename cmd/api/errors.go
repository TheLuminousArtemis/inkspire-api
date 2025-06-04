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
	writeJSONError(w, http.StatusBadGateway, err.Error())
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.l.Errorf("bad request error: %v path: %s err: %v", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) notFoundErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.l.Warnf("not found error: %v path: %s err: %v", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusNotFound, "not found error")
}

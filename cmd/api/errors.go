package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %v path: %s err: %v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "Internal server error")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %v path: %s err: %v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadGateway, err.Error())
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %v path: %s err: %v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) notFoundErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %v path: %s err: %v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, "not found error")
}

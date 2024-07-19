package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

var clientErrors = map[string]struct{}{
	"invalid input":                       {},
	"task already exists":                 {},
	"task doesn't exists or incorrect id": {},
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace) // Output() reports the file name and line number one step back in the stack trace
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal Server Error"))
}

func (app *application) clientError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(err.Error()))
}

func (app *application) handleError(w http.ResponseWriter, err error) {
	_, isClientError := clientErrors[err.Error()]
	if isClientError {
		app.clientError(w, err)
	} else {
		app.serverError(w, err)
	}
}

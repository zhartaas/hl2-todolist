package main

import (
	"github.com/swaggo/http-swagger" // http-swagger middleware)
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/create", app.createTask)
	mux.HandleFunc("/update", app.updateTask)
	mux.HandleFunc("/getTasks", app.getTasks)
	mux.HandleFunc("/delete", app.deleteTask)
	mux.HandleFunc("/taskDone", app.taskDone)

	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		httpSwagger.Handler(
			httpSwagger.URL("https://hl2-todolist.onrender.com/swagger/doc.json"),
		).ServeHTTP(w, r)
	})
	return mux
}

package main

import (
	"log"
	"net/http"
)

func (app *Application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
		"env":    app.Config.Addr,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		log.Print(err.Error())
		app.writeJsonError(w, http.StatusInternalServerError, err.Error())
	}
}

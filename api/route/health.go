package route

import (
	"log"
	"net/http"

	api_utils "github.com/KonovalovIly/anki_pdf/api/utils"
	"github.com/go-chi/chi/v5"
)

func (app *Application) HealthHandlersSetup(r chi.Router) {
	r.Get("/health", app.healthCheckHandler)
}

func (app *Application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
		"env":    app.Config.Addr,
	}

	if err := api_utils.JsonResponse(w, http.StatusOK, data); err != nil {
		log.Print(err.Error())
		api_utils.WriteJsonError(w, http.StatusInternalServerError, err)
	}
}

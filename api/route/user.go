package route

import (
	"fmt"
	"net/http"
	"strconv"

	api_model "github.com/KonovalovIly/anki_pdf/api/model"
	api_utils "github.com/KonovalovIly/anki_pdf/api/utils"
	"github.com/go-chi/chi/v5"
)

func (app *Application) UserHandlerSetup(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/register", app.registerUserHandler)

		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", app.userGetHandler)
		})
	})
}

func (app *Application) userGetHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil || userID <= 0 {
		api_utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	userDto, e := app.Storage.User.GetUser(r.Context(), userID)
	if e != nil {
		api_utils.WriteJsonDatabaseError(w, http.StatusInternalServerError, e)
		return
	}

	api_utils.JsonResponse(w, http.StatusAccepted, userDto)
}

func (app *Application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload api_model.UserRegisterPayload

	if err := api_utils.ReadJson(w, r, &payload); err != nil {
		api_utils.WriteJsonError(w, http.StatusBadRequest, err)
		return
	}

	if err := api_utils.Validator.Struct(payload); err != nil {
		api_utils.WriteJsonError(w, http.StatusBadRequest, err)
		return
	}
	userDto := payload.MapToDatabaseUser()

	e := app.Storage.User.SaveUser(r.Context(), userDto)
	if e != nil {
		api_utils.WriteJsonDatabaseError(w, http.StatusInternalServerError, e)
		return
	}
	api_utils.JsonResponse(w, http.StatusCreated, userDto)
}

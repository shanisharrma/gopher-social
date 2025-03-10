package main

import (
	"net/http"
)

// HealthCheck godoc
//
//	@Summary		Check health
//	@Description	Checks health of software
//	@Tags			Ops
//	@Produce		json
//	@Success		200	{string}	string	"OK"
//	@failure		500	{object}	error	"An error occured"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.jsonResponse(w, http.StatusOK, "Server is running fine", map[string]string{"version": version}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

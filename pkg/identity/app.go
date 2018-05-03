package identity

import (
	"log"
	"net/http"
)

type App struct {
	DB IdentityDatastore
}

func (app *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Printf("store.Get() failed with '%s'\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Bearer Token: %v\n", session.Values["BearerToken"])
}

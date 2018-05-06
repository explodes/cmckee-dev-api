package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cmckee-dev/cmckee-dev-api/pkg/identity"
	"github.com/gorilla/mux"
)

func main() {

	db, err := identity.InitDB(
		os.Getenv("PERSONA_DB_USER"),
		os.Getenv("PERSONA_DB_PASS"),
		os.Getenv("DEV_PERSONA_DB_DBNAME"))

	if err != nil {
		log.Fatal(err)
	}

	app := identity.App{db}

	// CR(explodes): consider a method that builds a router
	// func (a *App) CreateRouter() *mux.Router
	// so that this file doesn't need to perform wiring, app can manage its own
	// routing.
	//
	// side note: Another benefit to that is that you could use an interface:
	// type RouterCreator interface {
	// 		CreateRouter() *mux.Router
	// }
	// to compose several services together and serve them simultaneously, 
	// importing re-usable RouterCreators like a system-health service.
	//
	// That ideology is way of scope for this project, but it is a very useful pattern.

	router := mux.NewRouter()

	router.HandleFunc("/", app.HomeHandler)

	// Authentication

	router.HandleFunc("/login", app.LoginHandler)
	router.HandleFunc("/github_oauth_callback", app.GithubCallbackHandler)
	router.HandleFunc("/logout", app.LogoutHandler)

	// Users
	router.HandleFunc("/users", app.ReadAllUsersHandler).Methods("GET")
	router.HandleFunc("/user", app.CreateUserHandler).Methods("POST")
	router.HandleFunc("/user/{id}", app.ReadUserHandler).Methods("GET")
	router.HandleFunc("/user/{id}", app.UpdateUserHandler).Methods("POST")
	router.HandleFunc("/user/{id}", app.DeleteUserHandler).Methods("DELETE")

	log.Println("Starting server...")
	// CR(explodes): this pattern is often seen in the wild, and doesn't matter much here
	// without SIGINT handling, but unless the program is killed, it could log nil.
	// replace with:
	// if err := http.ListenAndServe(":4545", router) {

	}
	log.Fatal(http.ListenAndServe(":4545", router))
}

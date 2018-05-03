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
	log.Fatal(http.ListenAndServe(":4545", router))
}

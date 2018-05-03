package identity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

// DB Base Functionality

type User struct {
	ID        int       `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IdentityDatastore interface {
	getUsers(start, count int) ([]*User, error)
	createUser() error
	readUser() error
	updateUser() error
	deleteUser() error
}

type DB struct {
	*sql.DB
}

func InitDB(user, pass, dbname string) (*DB, error) {

	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, pass, dbname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// DB IO

func (db *DB) readAllUsers(start, count int) ([]*User, error) {
	rows, err := db.Query(
		"SELECT id, user_id, email, created_at, updated_at FROM users LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]*User, 0)

	for rows.Next() {
		usr := new(User)
		if err := rows.Scan(&usr.ID, &usr.UserID, &usr.Email, &usr.CreatedAt, &usr.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, usr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (db *DB) createUser(user *User) error {
	err := db.QueryRow(
		"INSERT INTO users(email) VALUES($1) RETURNING id",
		user.Email).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) readUser(user *User) error {
	err := db.QueryRow(
		"SELECT id, user_id, email, created_at, updated_at FROM users WHERE id=$1",
		user.ID).Scan(&user.ID, &user.UserID, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (db *DB) updateUser(user *User) error {
	_, err := db.Exec("UPDATE users SET email=$1, updated_at=$2",
		user.Email, time.Now())

	return err
}

func (db *DB) deleteUser(user *User) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", user.ID)

	return err
}

// Routing Handlers

func (app *App) ReadAllUsersHandler(w http.ResponseWriter, r *http.Request) {

	users, err := app.DB.getUsers(0, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (app *App) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	if err := app.DB.createUser(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (app *App) ReadUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid User Id", http.StatusBadRequest)
		return
	}

	user := &User{ID: id}

	if err := app.DB.readUser(user); err != nil {
		switch err {
		case sql.ErrNoRows:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (app *App) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
	return
}

func (app *App) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
	return
}

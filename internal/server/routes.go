package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var user = &User{}

func setupRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", rootHandler).Methods("GET").Name("root")
	router.HandleFunc("/users", getUsersHandler(db)).Methods("GET").Name("users")
	router.HandleFunc("/adduser", addUserHandler(db)).Methods("POST").Name("adduser")
	router.HandleFunc("/deleteuser", deleteUserHandler(db)).Methods("POST").Name("deleteuser")
	router.HandleFunc("/login", loginHandler(db)).Methods("POST").Name("login")
	router.HandleFunc("/api/v1/healthcheck", healthCheckHandler).Methods("GET").Name("healthcheck")

	private := router.PathPrefix("/private").Subrouter()
	private.Use(authMiddleware)
	private.HandleFunc("/admin", adminHandler).Methods("GET").Name("admin")

	return router
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Welcome to CloudFiles API`))
}

func getUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT name, password FROM users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		users := make([]map[string]interface{}, 0, 16)
		for rows.Next() {
			if err := rows.Scan(user.Userame, user.Password); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			user := map[string]interface{}{
				"name":     user.Userame,
				"password": user.Password,
			}
			users = append(users, user)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func addUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		password := r.FormValue("password")

		if name == "" || password == "" {
			http.Error(w, "Name and password are required", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", name, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func deleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.FormValue("id")
		if idStr == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

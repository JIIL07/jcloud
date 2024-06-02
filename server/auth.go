package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var data auth

func (app *application) UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)

	case http.MethodPost:

		var dataP auth

		err := json.NewDecoder(r.Body).Decode(&dataP)
		if err != nil {
			http.Error(w, "Can't read data from request", http.StatusBadRequest)
			return
		}

		// Сохранение данных
		data = dataP

		// Отправка сообщения об успешном сохранении
		fmt.Fprintf(w, "Data saved successfully")

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

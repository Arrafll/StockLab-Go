package services

import (
	"encoding/json"
	"net/http"
)

func GetUserList(w http.ResponseWriter, r *http.Request) {
	users := map[string]string{
		"User":  "Andre",
		"User2": "Rafli",
		"User3": "Brother",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

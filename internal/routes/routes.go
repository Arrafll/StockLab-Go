package routes

import (
	"encoding/json"
	"net/http"

	userService "github.com/Arrafll/setoko-go.git/internal/services/user"
)

func RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("/user", userService.GetUserList)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Endpoint not found",
		})
	})
}

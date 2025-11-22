package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, data any) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("encode response:", err)
	}
}

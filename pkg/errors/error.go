package errors

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Data string `json:"message"`
}

func Response(w http.ResponseWriter, err error, statusCode int) {
	output := &Error{Data: err.Error()}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(output)
}

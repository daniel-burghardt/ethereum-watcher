package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/daniel-burghardt/ethereum-parser/data"
)

type Handler struct {
	Repo data.Repository
}

func InternalServerError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusBadRequest)
}

func Created(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

func Json(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return fmt.Errorf("could not encode json data: %w", err)
	}

	return nil
}

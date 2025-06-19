package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ReadJSON(w http.ResponseWriter, r *http.Request, out any) error {
	const op = "utils.ReadJSON"

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("op: %s, err: %w", op, err)
	}
	defer r.Body.Close()

	err = json.Unmarshal(data, out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("op: %s, err: %w", op, err)
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, data any, status int) error {
	const op = "utils.WriteJSON"

	out, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("op: %s, err: %w", op, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)

	return nil
}

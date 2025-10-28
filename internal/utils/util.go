package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, data any) error {
	dat, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = w.Write(dat)
	if err != nil {
		return err
	}
	return nil
}

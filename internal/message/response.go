package message

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JsonMessage struct {
	ID int `json:"id"`
}

func MakeResponse(writer http.ResponseWriter, ID int, statusCode int) error {
	msg := JsonMessage{ID: ID}

	out, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("Couldn't marshall. error message: %s", err)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	writer.Write(out)

	return nil
}

package types

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	PodName      string            `json:"name"`
	PodNamespace string            `json:"namespace"`
	ID           string            `json:"id"`
	Args         map[string]string `json:"args"`
}

func FromRequest(r *http.Request) (Request, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return Request{}, fmt.Errorf("reading request body from %s: %w", r.RemoteAddr, err)
	}
	defer r.Body.Close()

	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		return Request{}, fmt.Errorf("unmarshalling request body from %s: %w", r.RemoteAddr, err)
	}

	return req, nil
}

func (r Request) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

func (r Request) Validate() error {
	if r.PodName == "" {
		return fmt.Errorf("missing pod name")
	}
	if r.PodNamespace == "" {
		return fmt.Errorf("missing pod namespace")
	}
	if r.ID == "" {
		return fmt.Errorf("missing id")
	}
	return nil
}

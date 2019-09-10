package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) handleError(err error, code int) http.HandlerFunc {
	type errRes struct {
		Code  int    `json:"statusCode"`
		Error string `json:"error"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		cstErr := &errRes{
			Code:  code,
			Error: err.Error(),
		}

		fmt.Printf("[%d] error: %s\n", code, err.Error())

		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(&cstErr); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

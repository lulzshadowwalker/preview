package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	preview "github.com/lulzshadowwalker/preview/pkg"
)

type ApiError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e ApiError) Error() string {
	return e.Message
}

type ErronousHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func unwrap(fn ErronousHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
      slog.Error("failed to get preview", "err", err)

			if apiErr, ok := err.(ApiError); ok {
				http.Error(w, apiErr.Message, apiErr.Status)
				return
			}

			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

const port = 8712

func main() {
  defer preview.Close()

	http.HandleFunc("GET /preview", unwrap(handleGetPreview))

	slog.Info("server started", "port", port, "url", fmt.Sprintf("http://localhost:%d", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			slog.Info("server shutdown")
			return
		}

		slog.Error("server crashed", "err", err)
	}
}

func handleGetPreview(w http.ResponseWriter, r *http.Request) error {
  u := r.URL.Query().Get("url")
	if _, err := url.Parse(u); err != nil {
		return ApiError{
			Message: "invalid url",
			Status:  http.StatusBadRequest,
		}
	}

  target, err := url.QueryUnescape(u)
  if err != nil {
    return err
  }

	p, err := preview.FromURL(target)
	if err != nil {
		if errors.Is(err, preview.ErrNotFound) {
			return ApiError{
				Message: "not found",
				Status:  http.StatusNotFound,
			}
		}

		return err
	}

  w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(map[string]any{
		"data": map[string]any{
			"preview": p,
		},
	})
}

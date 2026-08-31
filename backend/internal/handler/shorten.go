package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/DragonZero000/url-shortener/internal/service"
)

type URLRequest struct {
	URL string `json:"url"`
}

type ShortedURL struct {
	URL string `json:"url"`
}

type Handler struct {
	service    *service.Service
	webBaseURL string
}

func NewHandler(s *service.Service, webBaseURL string) *Handler {
	return &Handler{service: s, webBaseURL: webBaseURL}
}

func (h *Handler) PostShorten(w http.ResponseWriter, r *http.Request) {
	var URL URLRequest
	err := json.NewDecoder(r.Body).Decode(&URL)
	if err != nil {
		slog.Warn("invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	code, err := h.service.Shorten(URL.URL)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) {
			slog.Warn("invalid URL", "error", err)
			http.Error(w, "invalid URL error", http.StatusBadRequest)
			return
		}
		slog.Error("failed to generate short URL", "url", URL.URL, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	shortenURL := h.webBaseURL + "/" + code
	ShortedURL := ShortedURL{URL: shortenURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ShortedURL)
}

func (h *Handler) GetByShorten(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	originalURL, err := h.service.GetByCode(code)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			slog.Warn("the link you were looking for was not found", "error", err)
			http.Error(w, "the link you were looking for was not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to resolve short code", "code", code, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusSeeOther)
}

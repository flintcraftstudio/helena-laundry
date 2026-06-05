package handler

import (
	"log/slog"
	"net/http"

	"github.com/firefly-software-mt/advanced-template/internal/view"
)

// Services handles GET /services and renders the services & pricing page.
func Services() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := view.ServicesPage().Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// About handles GET /about and renders Chanté's story.
func About() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := view.AboutPage().Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// Businesses handles GET /businesses and renders the B2B page.
func Businesses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := view.BusinessesPage().Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// NotFound handles unmatched GET paths with a branded 404 page.
func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if err := view.NotFoundPage().Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

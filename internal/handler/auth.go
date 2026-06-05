package handler

import (
	"log/slog"
	"net/http"

	"github.com/firefly-software-mt/advanced-template/internal/session"
	"github.com/firefly-software-mt/advanced-template/internal/store"
	"github.com/firefly-software-mt/advanced-template/internal/view"

	"golang.org/x/crypto/bcrypt"
)

// LoginPage handles GET /login and renders the login form.
func LoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if session.FromContext(r.Context()) != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if err := view.LoginPage("", "").Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// LoginSubmit handles POST /login, validates credentials, creates a session, and redirects.
func LoginSubmit(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			if err := view.LoginForm("Email and password are required.", email).Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}

		userID, _, passwordHash, err := s.GetUserByEmail(r.Context(), email)
		if err != nil {
			if err := view.LoginForm("Invalid email or password.", email).Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
			if err := view.LoginForm("Invalid email or password.", email).Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}

		if err := session.Create(r.Context(), w, s, userID); err != nil {
			slog.Error("session create error", "err", err)
			if err := view.LoginForm("Something went wrong. Please try again.", email).Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}

		redirect(w, r, "/")
	}
}

// Logout handles POST /logout, destroys the session, and redirects.
func Logout(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := session.Destroy(r.Context(), w, r, s); err != nil {
			slog.Error("session destroy error", "err", err)
		}
		redirect(w, r, "/login")
	}
}

// redirect sends the client to url. For htmx requests it must use the HX-Redirect
// header on a 2xx response — an http.Redirect (3xx) is followed transparently by
// the XHR, so htmx never sees it and swaps the destination page into the form.
// Non-htmx requests get a normal 303 redirect.
func redirect(w http.ResponseWriter, r *http.Request, url string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", url)
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
}

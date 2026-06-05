package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/firefly-software-mt/advanced-template/internal/mail"
	"github.com/firefly-software-mt/advanced-template/internal/store"
	"github.com/firefly-software-mt/advanced-template/internal/view"
)

// Contact handles GET /contact and renders the contact / schedule page.
func Contact() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := view.ContactPage(nil, nil, false, nil, nil, false).Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// QuestionSubmit handles POST /contact/question — the question / feedback form.
// Stores the submission and (if configured) emails Chanté; replies with the
// htmx-swapped confirmation or the form with inline errors.
func QuestionSubmit(st *store.Store, mailer *mail.Client, turnstileSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		values := map[string]string{
			"name":    strings.TrimSpace(r.FormValue("name")),
			"contact": strings.TrimSpace(r.FormValue("contact")),
			"topic":   strings.TrimSpace(r.FormValue("topic")),
			"message": strings.TrimSpace(r.FormValue("message")),
		}

		errors := map[string]string{}
		if values["name"] == "" {
			errors["name"] = "I'll need a name to know who I'm writing back."
		}
		if values["contact"] == "" {
			errors["contact"] = "Leave a phone or email so I can reach you."
		}
		if values["message"] == "" {
			errors["message"] = "Tell me what's on your mind."
		}
		if values["topic"] == "" {
			values["topic"] = "question"
		}

		if len(errors) > 0 {
			renderQuestionForm(w, r, errors, values, false)
			return
		}

		if turnstileSecret != "" && !verifyTurnstile(turnstileSecret, r.FormValue("cf-turnstile-response"), r.RemoteAddr) {
			renderQuestionForm(w, r, map[string]string{"form": "That verification didn't take. Give it another try."}, values, false)
			return
		}

		if err := st.CreateQuestion(r.Context(), values["name"], values["contact"], values["topic"], values["message"]); err != nil {
			slog.Error("store question error", "err", err)
			renderQuestionForm(w, r, map[string]string{"form": "Something went sideways on my end. Try again, or just call or text me."}, values, false)
			return
		}

		if mailer != nil {
			body := fmt.Sprintf("Topic: %s\nReach them at: %s\n\n%s", values["topic"], values["contact"], values["message"])
			msg := mail.Message{
				Name:    values["name"],
				Email:   emailIfPresent(values["contact"]),
				Subject: fmt.Sprintf("Helena Laundry — %s from %s", values["topic"], values["name"]),
				Body:    body,
			}
			if err := mailer.Send(msg); err != nil {
				slog.Error("postmark send error (question)", "err", err)
				// Submission is already stored; don't fail the visitor over email.
			}
		}

		renderQuestionForm(w, r, nil, nil, true)
	}
}

// WaitlistSubmit handles POST /contact/waitlist — the out-of-area waitlist form.
func WaitlistSubmit(st *store.Store, mailer *mail.Client, turnstileSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		values := map[string]string{
			"name":     strings.TrimSpace(r.FormValue("name")),
			"location": strings.TrimSpace(r.FormValue("location")),
			"contact":  strings.TrimSpace(r.FormValue("contact")),
			"kind":     strings.TrimSpace(r.FormValue("kind")),
			"notes":    strings.TrimSpace(r.FormValue("notes")),
		}

		errors := map[string]string{}
		if values["name"] == "" {
			errors["name"] = "I'll need a name for the list."
		}
		if values["location"] == "" {
			errors["location"] = "Tell me where you are — that's how I decide where to head next."
		}
		if values["contact"] == "" {
			errors["contact"] = "Leave a phone or email so I can reach you when I'm headed your way."
		}
		if values["kind"] != "business" {
			values["kind"] = "residential"
		}

		if len(errors) > 0 {
			renderWaitlistForm(w, r, errors, values, false)
			return
		}

		if turnstileSecret != "" && !verifyTurnstile(turnstileSecret, r.FormValue("cf-turnstile-response"), r.RemoteAddr) {
			renderWaitlistForm(w, r, map[string]string{"form": "That verification didn't take. Give it another try."}, values, false)
			return
		}

		if err := st.CreateWaitlistEntry(r.Context(), values["name"], values["location"], values["contact"], values["kind"], values["notes"]); err != nil {
			slog.Error("store waitlist error", "err", err)
			renderWaitlistForm(w, r, map[string]string{"form": "Something went sideways on my end. Try again, or just call or text me."}, values, false)
			return
		}

		if mailer != nil {
			body := fmt.Sprintf("Location: %s\nType: %s\nReach them at: %s\n\nNotes: %s", values["location"], values["kind"], values["contact"], values["notes"])
			msg := mail.Message{
				Name:    values["name"],
				Email:   emailIfPresent(values["contact"]),
				Subject: fmt.Sprintf("Helena Laundry waitlist — %s (%s)", values["location"], values["kind"]),
				Body:    body,
			}
			if err := mailer.Send(msg); err != nil {
				slog.Error("postmark send error (waitlist)", "err", err)
			}
		}

		renderWaitlistForm(w, r, nil, nil, true)
	}
}

func renderQuestionForm(w http.ResponseWriter, r *http.Request, errs, values map[string]string, success bool) {
	if err := view.QuestionForm(errs, values, success).Render(r.Context(), w); err != nil {
		slog.Error("render error", "err", err)
	}
}

func renderWaitlistForm(w http.ResponseWriter, r *http.Request, errs, values map[string]string, success bool) {
	if err := view.WaitlistForm(errs, values, success).Render(r.Context(), w); err != nil {
		slog.Error("render error", "err", err)
	}
}

// emailIfPresent returns the contact string as a ReplyTo only when it looks like
// an email; Postmark rejects a non-email ReplyTo. Phone-only contacts get "".
func emailIfPresent(contact string) string {
	if strings.Contains(contact, "@") {
		return contact
	}
	return ""
}

// verifyTurnstile checks a Turnstile token against the Cloudflare API.
func verifyTurnstile(secret, token, remoteIP string) bool {
	resp, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", url.Values{
		"secret":   {secret},
		"response": {token},
		"remoteip": {remoteIP},
	})
	if err != nil {
		slog.Error("turnstile verify request failed", "err", err)
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("turnstile verify decode failed", "err", err)
		return false
	}
	if !result.Success {
		slog.Warn("turnstile verification failed")
	}
	return result.Success
}

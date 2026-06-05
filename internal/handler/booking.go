package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/firefly-software-mt/advanced-template/internal/mail"
	"github.com/firefly-software-mt/advanced-template/internal/store"
)

// validPickupDays are the weekday options offered in the booking modal.
var validPickupDays = map[string]bool{
	"Monday": true, "Tuesday": true, "Wednesday": true, "Thursday": true, "Friday": true,
}

// planLabels maps a booking plan key to the human label used in the email.
var planLabels = map[string]string{
	"standard": "Wash, dry, fold & hang — $2.00/lb (one-time)",
	"weekly":   "Standing weekly pickup — $2.00/lb",
	"rush":     "Rush — $3.00/lb (back in 24 hours)",
}

// BookingSubmit handles POST /contact/booking — the multi-step "Book a pickup"
// modal. The modal drives its own steps and confirmation client-side, so this
// endpoint only validates, persists, and notifies. Three outcomes:
//   - 204 No Content                  → booked; modal shows the confirmation.
//   - 200 + X-Booking-Result:out-of-area → ZIP outside the service area; nothing
//     stored, modal refers them to the waitlist.
//   - 4xx/5xx + plain message         → modal shows it as an inline banner.
//
// allowedZips is the set of ZIPs Chanté picks up from (from config).
func BookingSubmit(st *store.Store, mailer *mail.Client, turnstileSecret string, allowedZips map[string]bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			bookingError(w, http.StatusBadRequest, "That didn't come through. Give it another try.")
			return
		}

		name := strings.TrimSpace(r.FormValue("name"))
		contact := strings.TrimSpace(r.FormValue("contact"))
		day := strings.TrimSpace(r.FormValue("day"))
		plan := strings.TrimSpace(r.FormValue("plan"))
		address := strings.TrimSpace(r.FormValue("address"))
		zip := normalizeZip(r.FormValue("zip"))
		notes := strings.TrimSpace(r.FormValue("notes"))

		// The modal validates required fields client-side, so a miss here means
		// JS was bypassed — answer plainly rather than with field-level errors.
		if name == "" || contact == "" || address == "" {
			bookingError(w, http.StatusUnprocessableEntity, "I'll need a name, a way to reach you, and a pickup address. Add those and try again.")
			return
		}
		if zip == "" {
			bookingError(w, http.StatusUnprocessableEntity, "I'll need a 5-digit ZIP so I know if I reach you.")
			return
		}
		if !validPickupDays[day] {
			day = "Tuesday"
		}
		if _, ok := planLabels[plan]; !ok {
			plan = "standard"
		}

		if turnstileSecret != "" && !verifyTurnstile(turnstileSecret, r.FormValue("cf-turnstile-response"), r.RemoteAddr) {
			bookingError(w, http.StatusUnprocessableEntity, "That verification didn't take. Give it another try, or just call or text me.")
			return
		}

		// Service-area gate: a ZIP outside the route isn't a booking — point them
		// at the waitlist instead. Nothing is stored; the waitlist captures intent.
		if !allowedZips[zip] {
			slog.Info("booking out of area", "zip", zip)
			w.Header().Set("X-Booking-Result", "out-of-area")
			w.WriteHeader(http.StatusOK)
			return
		}

		if err := st.CreateBooking(r.Context(), name, contact, day, plan, address, zip, notes); err != nil {
			slog.Error("store booking error", "err", err)
			bookingError(w, http.StatusInternalServerError, "Something went sideways on my end. Try again, or just call or text me.")
			return
		}

		if mailer != nil {
			body := fmt.Sprintf(
				"New pickup booking.\n\nName: %s\nReach them at: %s\nPickup day: %s\nPlan: %s\nAddress: %s\nZIP: %s\n\nNotes: %s",
				name, contact, day, planLabels[plan], orDash(address), zip, orDash(notes),
			)
			msg := mail.Message{
				Name:    name,
				Email:   emailIfPresent(contact),
				Subject: fmt.Sprintf("Helena Laundry — pickup booked: %s (%s)", name, day),
				Body:    body,
			}
			if err := mailer.Send(msg); err != nil {
				slog.Error("postmark send error (booking)", "err", err)
				// Booking is already stored; don't fail the visitor over email.
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// normalizeZip extracts the leading 5-digit ZIP from raw input, tolerating
// surrounding text and ZIP+4 ("59601-1234" -> "59601"). Returns "" if there
// aren't 5 leading digits.
func normalizeZip(raw string) string {
	var digits strings.Builder
	for _, r := range raw {
		if r >= '0' && r <= '9' {
			digits.WriteRune(r)
			if digits.Len() == 5 {
				break
			}
		} else if digits.Len() > 0 {
			// Stop at the first non-digit after digits started (drops the +4).
			break
		}
	}
	if digits.Len() != 5 {
		return ""
	}
	return digits.String()
}

// bookingError writes a plain-text message the booking modal surfaces inline.
func bookingError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	_, _ = w.Write([]byte(msg))
}

// orDash returns "—" for an empty string so emails read cleanly.
func orDash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

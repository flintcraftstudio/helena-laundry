package handler

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/firefly-software-mt/advanced-template/internal/settings"
	"github.com/firefly-software-mt/advanced-template/internal/store"
	"github.com/firefly-software-mt/advanced-template/internal/view"
)

// AdminSettings handles GET /admin/settings — the site settings editor.
func AdminSettings(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := st.GetSettings(r.Context())
		if err != nil {
			adminServerError(w, r, "load settings", err)
			return
		}
		render(w, r, view.AdminSettingsPage(userEmail(r), s))
	}
}

// AdminSettingsUpdate handles POST /admin/settings. It validates first; only a
// clean form is written, the live cache + view vars refreshed, and the success
// toast fired. A form with bad input is re-rendered with the offending fields
// highlighted — nothing reaches the public site until it's valid.
func AdminSettingsUpdate(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		s := store.Settings{
			SiteName:          strings.TrimSpace(r.FormValue("site_name")),
			Tagline:           strings.TrimSpace(r.FormValue("tagline")),
			Phone:             strings.TrimSpace(r.FormValue("phone")),
			Email:             strings.TrimSpace(r.FormValue("email")),
			Hours:             strings.TrimSpace(r.FormValue("hours")),
			Turnaround:        strings.TrimSpace(r.FormValue("turnaround")),
			AcceptingBookings: r.FormValue("accepting_bookings") != "",
			LbRate:            strings.TrimSpace(r.FormValue("lb_rate")),
			RushRate:          strings.TrimSpace(r.FormValue("rush_rate")),
			MinimumOrder:      strings.TrimSpace(r.FormValue("minimum_order")),
			AllowedZips:       parseZipList(r.FormValue("allowed_zips")),
			PickupDays:        r.Form["pickup_days"],
		}
		rates := pairFlatRates(r.Form["flat_label"], r.Form["flat_price"])

		// Validate + normalize in place. On any problem, hand the submitted values
		// straight back so nothing is lost and nothing goes live.
		if errs := validateSettings(&s, rates); len(errs) > 0 {
			s.FlatRates = rates
			w.Header().Set("HX-Trigger", toastSettingsError)
			render(w, r, view.SettingsRegion(s, errs))
			return
		}

		if err := st.UpdateSettings(r.Context(), s); err != nil {
			adminServerError(w, r, "save settings", err)
			return
		}
		if err := st.ReplaceFlatRates(r.Context(), rates); err != nil {
			adminServerError(w, r, "save flat rates", err)
			return
		}

		// Refresh the live cache + the public-facing view vars.
		full, err := st.GetSettings(r.Context())
		if err != nil {
			adminServerError(w, r, "reload settings", err)
			return
		}
		settings.Set(full)
		view.ApplySettings(full)

		w.Header().Set("HX-Trigger", toastSaved)
		render(w, r, view.SettingsRegion(full, nil))
	}
}

// AdminSettingsFlatRateRow handles GET /admin/settings/flat-rate-row — returns a
// single blank flat-rate row so "Add item" can append one without a full save.
func AdminSettingsFlatRateRow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, r, view.SettingsFlatRateRow(store.FlatRate{}))
	}
}

var settingsEmailRe = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)

// validateSettings normalizes the values that flow to the live site and returns
// a field->message map of anything it couldn't accept. Money fields are cleaned
// in place (a stray "$" or spaces are forgiven, not rejected); blanks are always
// allowed so the operator can leave a price unset.
func validateSettings(s *store.Settings, rates []store.FlatRate) map[string]string {
	errs := map[string]string{}

	if v, ok := normalizeMoney(s.LbRate); ok {
		s.LbRate = v
	} else {
		errs["lb_rate"] = "Use a number, like 2.50."
	}
	if v, ok := normalizeMoney(s.RushRate); ok {
		s.RushRate = v
	} else {
		errs["rush_rate"] = "Use a number, like 1.00."
	}
	if v, ok := normalizeMoney(s.MinimumOrder); ok {
		s.MinimumOrder = v
	} else {
		errs["minimum_order"] = "Use a number, like 25."
	}

	if s.Email != "" && !settingsEmailRe.MatchString(s.Email) {
		errs["email"] = "That doesn't look like an email address."
	}
	if s.Phone != "" && countDigits(s.Phone) < 10 {
		errs["phone"] = "Add a full phone number, area code and all."
	}

	for i := range rates {
		if strings.TrimSpace(rates[i].Label) == "" {
			continue // blank rows are dropped on save
		}
		if v, ok := normalizeMoney(rates[i].Price); ok {
			rates[i].Price = v
		} else {
			errs["flat_rates"] = "One of the flat-rate prices isn't a number."
		}
	}

	return errs
}

// normalizeMoney accepts a blank, or a positive dollar amount with an optional
// leading "$" and surrounding spaces, returning the bare numeric string. ok is
// false only when there's non-numeric junk we shouldn't quietly publish.
func normalizeMoney(raw string) (string, bool) {
	v := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(raw), "$"))
	v = strings.TrimSpace(v)
	if v == "" {
		return "", true
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || f < 0 {
		return raw, false
	}
	return v, true
}

// countDigits counts the 0-9 runes in s (used for a light phone sanity check).
func countDigits(s string) int {
	n := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			n++
		}
	}
	return n
}

// parseZipList splits a comma/space/newline-separated list, keeping only valid
// 5-digit ZIPs so a stray character can't quietly break the service area.
func parseZipList(raw string) []string {
	tokens := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\n' || r == '\t' || r == ';'
	})
	var out []string
	for _, t := range tokens {
		if z := normalizeZip(t); z != "" {
			out = append(out, z)
		}
	}
	return out
}

// pairFlatRates zips the parallel label/price form arrays into flat-rate rows.
func pairFlatRates(labels, prices []string) []store.FlatRate {
	out := make([]store.FlatRate, 0, len(labels))
	for i, label := range labels {
		price := ""
		if i < len(prices) {
			price = prices[i]
		}
		out = append(out, store.FlatRate{Label: label, Price: price})
	}
	return out
}

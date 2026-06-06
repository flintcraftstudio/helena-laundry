package handler

import (
	"net/http"
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

// AdminSettingsUpdate handles POST /admin/settings — save, refresh the live
// cache + view vars, and re-render the form with a toast.
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

		if err := st.UpdateSettings(r.Context(), s); err != nil {
			adminServerError(w, r, "save settings", err)
			return
		}
		if err := st.ReplaceFlatRates(r.Context(), pairFlatRates(r.Form["flat_label"], r.Form["flat_price"])); err != nil {
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
		render(w, r, view.SettingsRegion(full))
	}
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

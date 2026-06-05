package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/firefly-software-mt/advanced-template/internal/store"
	"github.com/firefly-software-mt/advanced-template/internal/view"
)

// buildCalendar fetches the scheduled + unscheduled bookings for a month and
// assembles the view model. Shared by the page, htmx nav, and post-update render.
func buildCalendar(st *store.Store, r *http.Request, month string) (view.CalMonth, error) {
	now := time.Now()
	start, end := view.CalendarQueryRange(month, now)
	scheduled, err := st.ListBookingsScheduledBetween(r.Context(), start, end)
	if err != nil {
		return view.CalMonth{}, err
	}
	unscheduled, err := st.ListUnscheduledBookings(r.Context())
	if err != nil {
		return view.CalMonth{}, err
	}
	return view.BuildCalendar(month, scheduled, unscheduled, now), nil
}

// AdminCalendar handles GET /admin/calendar?m=YYYY-MM. htmx (month nav) gets just
// the calendar region; a full navigation gets the whole page.
func AdminCalendar(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cal, err := buildCalendar(st, r, r.URL.Query().Get("m"))
		if err != nil {
			adminServerError(w, r, "build calendar", err)
			return
		}
		if isHTMX(r) {
			render(w, r, view.CalendarRegion(cal))
			return
		}
		render(w, r, view.AdminCalendarPage(userEmail(r), cal))
	}
}

// AdminCalendarPickup handles GET /admin/calendar/pickup/{id}?m=YYYY-MM — the
// reschedule/status editor panel for one booking.
func AdminCalendarPickup(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, ok := getBooking(w, r, st)
		if !ok {
			return
		}
		render(w, r, view.CalendarPickupPanel(b, r.URL.Query().Get("m")))
	}
}

// AdminCalendarPickupUpdate handles POST /admin/calendar/pickup/{id} — save the
// date/window/status/notes, then re-render the calendar region for month "m".
func AdminCalendarPickupUpdate(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := adminID(w, r)
		if !ok {
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		status := validStatus(r.FormValue("status"), bookingStatuses)
		if status == "" {
			status = "new"
		}
		err := st.UpdateBooking(r.Context(), id,
			status,
			strings.TrimSpace(r.FormValue("scheduled_date")),
			strings.TrimSpace(r.FormValue("pickup_window")),
			strings.TrimSpace(r.FormValue("admin_notes")),
		)
		if err != nil {
			adminServerError(w, r, "update pickup", err)
			return
		}
		cal, err := buildCalendar(st, r, r.FormValue("m"))
		if err != nil {
			adminServerError(w, r, "rebuild calendar", err)
			return
		}
		w.Header().Set("HX-Trigger", toastSaved)
		render(w, r, view.CalendarRegion(cal))
	}
}

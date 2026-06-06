package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/firefly-software-mt/advanced-template/internal/store"
	"github.com/firefly-software-mt/advanced-template/internal/view"
)

// buildCalendar fetches the scheduled + unscheduled bookings for a view+anchor
// and assembles the view model. Shared by the page, htmx nav, and post-update
// render. vmode is "day"|"week"|"month" (defaulted by the view layer); anchor is
// a YYYY-MM-DD date the view centers on.
func buildCalendar(st *store.Store, r *http.Request, vmode, anchor string) (view.Cal, error) {
	now := time.Now()
	start, end := view.CalRange(vmode, anchor, now)
	scheduled, err := st.ListBookingsScheduledBetween(r.Context(), start, end)
	if err != nil {
		return view.Cal{}, err
	}
	unscheduled, err := st.ListUnscheduledBookings(r.Context())
	if err != nil {
		return view.Cal{}, err
	}
	return view.BuildCal(vmode, anchor, scheduled, unscheduled, now), nil
}

// AdminCalendar handles GET /admin/calendar?v=week&d=YYYY-MM-DD. htmx (view/date
// nav) gets just the calendar region; a full navigation gets the whole page.
func AdminCalendar(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		cal, err := buildCalendar(st, r, q.Get("v"), q.Get("d"))
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

// AdminCalendarDay handles GET /admin/calendar/day/{date}?v=&d= — the full pickup
// list for one day, loaded into the drawer by a "+N more" tap on a packed month
// cell. The v/d query carries the calendar scope to restore after an edit.
func AdminCalendarDay(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date := r.PathValue("date")
		day, err := time.Parse("2006-01-02", date)
		if err != nil {
			http.Error(w, "bad date", http.StatusBadRequest)
			return
		}
		next := day.AddDate(0, 0, 1).Format("2006-01-02")
		bookings, err := st.ListBookingsScheduledBetween(r.Context(), date, next)
		if err != nil {
			adminServerError(w, r, "list day pickups", err)
			return
		}
		render(w, r, view.CalendarDayPanel(view.DayPanelLabel(date), bookings, calRef(r.URL.Query().Get("v"), r.URL.Query().Get("d"))))
	}
}

// AdminCalendarPickup handles GET /admin/calendar/pickup/{id}?v=&d= — the
// reschedule/status editor panel for one booking.
func AdminCalendarPickup(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, ok := getBooking(w, r, st)
		if !ok {
			return
		}
		render(w, r, view.CalendarPickupPanel(b, calRef(r.URL.Query().Get("v"), r.URL.Query().Get("d"))))
	}
}

// AdminCalendarPickupUpdate handles POST /admin/calendar/pickup/{id} — save the
// date/window/status/notes, then re-render the calendar region for the scope the
// operator was on (v/d).
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
		cal, err := buildCalendar(st, r, r.FormValue("v"), r.FormValue("d"))
		if err != nil {
			adminServerError(w, r, "rebuild calendar", err)
			return
		}
		w.Header().Set("HX-Trigger", toastSavedClose)
		render(w, r, view.CalendarRegion(cal))
	}
}

// calRef builds the view-layer scope from the raw v/d request params.
func calRef(vmode, anchor string) view.CalRef {
	return view.CalRef{View: vmode, Anchor: anchor}
}

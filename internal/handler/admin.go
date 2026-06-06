package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"

	"github.com/firefly-software-mt/advanced-template/internal/session"
	"github.com/firefly-software-mt/advanced-template/internal/store"
	"github.com/firefly-software-mt/advanced-template/internal/view"
)

const (
	adminListLimit   = 200
	adminRecentLimit = 5
	adminPageSize    = 25
	toastSaved       = `{"flint:toast":{"variant":"success","title":"Saved","body":"Your site's updated."}}`
	// toastSavedClose also slides the edit drawer shut (used by the row updates).
	toastSavedClose = `{"flint:toast":{"variant":"success","title":"Saved"},"admin-drawer-close":true}`
	// toastSettingsError fires when a save is rejected for bad input — the form
	// comes back with the offending fields highlighted instead of going live.
	toastSettingsError = `{"flint:toast":{"variant":"danger","title":"Not saved","body":"Check the highlighted fields."}}`
)

// pageParam reads the 1-based ?page query value, defaulting to 1.
func pageParam(r *http.Request) int {
	if n, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && n > 1 {
		return n
	}
	return 1
}

// allowed status values per resource (mirrors the view vocabularies).
var (
	bookingStatuses  = set("new", "scheduled", "picked_up", "in_progress", "ready", "delivered", "canceled")
	questionStatuses = set("new", "handled", "archived")
	waitlistStatuses = set("new", "contacted", "archived")
)

func set(vals ...string) map[string]bool {
	m := make(map[string]bool, len(vals))
	for _, v := range vals {
		m[v] = true
	}
	return m
}

// AdminDashboard handles GET /admin — the operator overview.
func AdminDashboard(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		counts := view.AdminCounts{
			Bookings:  countOrZero(st, ctx, "bookings", "new"),
			Inquiries: countOrZero(st, ctx, "questions", "new"),
			Waitlist:  countOrZero(st, ctx, "waitlist", "new"),
		}
		bookings, _ := st.ListBookings(ctx, "", adminRecentLimit)
		inquiries, _ := st.ListQuestions(ctx, "", adminRecentLimit)
		waitlist, _ := st.ListWaitlist(ctx, "", adminRecentLimit)
		render(w, r, view.AdminDashboardPage(userEmail(r), counts, bookings, inquiries, waitlist))
	}
}

// --- Bookings ---------------------------------------------------------------

// AdminBookings handles GET /admin/bookings. htmx (tab) requests get just the
// section; a full navigation gets the whole page.
func AdminBookings(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := validStatus(r.URL.Query().Get("status"), bookingStatuses)
		search := strings.TrimSpace(r.URL.Query().Get("q"))
		page := pageParam(r)
		bookings, total, err := st.ListBookingsFiltered(r.Context(), store.BookingFilter{
			Status: status, Search: search,
			Limit: adminPageSize, Offset: (page - 1) * adminPageSize,
		})
		if err != nil {
			adminServerError(w, r, "load bookings", err)
			return
		}
		bv := view.BookingsView{
			Status: status, Search: search, Page: page,
			PageSize: adminPageSize, Total: total, Bookings: bookings,
		}
		if isHTMX(r) {
			render(w, r, view.BookingsSection(bv))
			return
		}
		render(w, r, view.AdminBookingsPage(userEmail(r), bv))
	}
}

// AdminBookingRow handles GET /admin/bookings/{id} — the summary row (cancel/refresh).
func AdminBookingRow(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, ok := getBooking(w, r, st)
		if !ok {
			return
		}
		render(w, r, view.BookingSummaryRow(b))
	}
}

// AdminBookingEdit handles GET /admin/bookings/{id}/edit — the inline edit form.
func AdminBookingEdit(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, ok := getBooking(w, r, st)
		if !ok {
			return
		}
		render(w, r, view.BookingEditPanel(b))
	}
}

// AdminBookingUpdate handles POST /admin/bookings/{id} — save status/schedule/notes.
func AdminBookingUpdate(st *store.Store) http.HandlerFunc {
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
			adminServerError(w, r, "update booking", err)
			return
		}
		b, err := st.GetBooking(r.Context(), id)
		if err != nil {
			adminServerError(w, r, "reload booking", err)
			return
		}
		w.Header().Set("HX-Trigger", toastSavedClose)
		render(w, r, view.BookingSummaryRow(b))
	}
}

// --- Inquiries --------------------------------------------------------------

func AdminInquiries(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := validStatus(r.URL.Query().Get("status"), questionStatuses)
		search := strings.TrimSpace(r.URL.Query().Get("q"))
		page := pageParam(r)
		items, total, err := st.ListQuestionsFiltered(r.Context(), store.QuestionFilter{
			Status: status, Search: search,
			Limit: adminPageSize, Offset: (page - 1) * adminPageSize,
		})
		if err != nil {
			adminServerError(w, r, "load inquiries", err)
			return
		}
		iv := view.InquiriesView{
			Status: status, Search: search, Page: page,
			PageSize: adminPageSize, Total: total, Inquiries: items,
		}
		if isHTMX(r) {
			render(w, r, view.InquiriesSection(iv))
			return
		}
		render(w, r, view.AdminInquiriesPage(userEmail(r), iv))
	}
}

func AdminInquiryRow(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, ok := getQuestion(w, r, st)
		if !ok {
			return
		}
		render(w, r, view.InquirySummaryRow(q))
	}
}

func AdminInquiryEdit(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, ok := getQuestion(w, r, st)
		if !ok {
			return
		}
		render(w, r, view.InquiryEditPanel(q))
	}
}

func AdminInquiryUpdate(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := adminID(w, r)
		if !ok {
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		status := validStatus(r.FormValue("status"), questionStatuses)
		if status == "" {
			status = "new"
		}
		if err := st.UpdateQuestionStatus(r.Context(), id, status, strings.TrimSpace(r.FormValue("admin_notes"))); err != nil {
			adminServerError(w, r, "update inquiry", err)
			return
		}
		q, err := st.GetQuestion(r.Context(), id)
		if err != nil {
			adminServerError(w, r, "reload inquiry", err)
			return
		}
		w.Header().Set("HX-Trigger", toastSavedClose)
		render(w, r, view.InquirySummaryRow(q))
	}
}

// --- Waitlist ---------------------------------------------------------------

func AdminWaitlist(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := validStatus(r.URL.Query().Get("status"), waitlistStatuses)
		loc := strings.TrimSpace(r.URL.Query().Get("loc"))
		items, err := st.ListWaitlist(r.Context(), status, adminListLimit)
		if err != nil {
			adminServerError(w, r, "load waitlist", err)
			return
		}
		wv := view.NewWaitlistView(status, loc, items)
		if isHTMX(r) {
			render(w, r, view.WaitlistSection(wv))
			return
		}
		render(w, r, view.AdminWaitlistPage(userEmail(r), wv))
	}
}

func AdminWaitlistRow(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e, ok := getWaitlist(w, r, st)
		if !ok {
			return
		}
		render(w, r, view.WaitlistSummaryRow(e))
	}
}

func AdminWaitlistEdit(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e, ok := getWaitlist(w, r, st)
		if !ok {
			return
		}
		render(w, r, view.WaitlistEditPanel(e))
	}
}

func AdminWaitlistUpdate(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := adminID(w, r)
		if !ok {
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		status := validStatus(r.FormValue("status"), waitlistStatuses)
		if status == "" {
			status = "new"
		}
		if err := st.UpdateWaitlistStatus(r.Context(), id, status, strings.TrimSpace(r.FormValue("admin_notes"))); err != nil {
			adminServerError(w, r, "update waitlist", err)
			return
		}
		e, err := st.GetWaitlistEntry(r.Context(), id)
		if err != nil {
			adminServerError(w, r, "reload waitlist", err)
			return
		}
		w.Header().Set("HX-Trigger", toastSavedClose)
		render(w, r, view.WaitlistSummaryRow(e))
	}
}

// --- shared helpers ---------------------------------------------------------

func getBooking(w http.ResponseWriter, r *http.Request, st *store.Store) (store.Booking, bool) {
	id, ok := adminID(w, r)
	if !ok {
		return store.Booking{}, false
	}
	b, err := st.GetBooking(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return store.Booking{}, false
	}
	return b, true
}

func getQuestion(w http.ResponseWriter, r *http.Request, st *store.Store) (store.Question, bool) {
	id, ok := adminID(w, r)
	if !ok {
		return store.Question{}, false
	}
	q, err := st.GetQuestion(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return store.Question{}, false
	}
	return q, true
}

func getWaitlist(w http.ResponseWriter, r *http.Request, st *store.Store) (store.WaitlistEntry, bool) {
	id, ok := adminID(w, r)
	if !ok {
		return store.WaitlistEntry{}, false
	}
	e, err := st.GetWaitlistEntry(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return store.WaitlistEntry{}, false
	}
	return e, true
}

// adminID parses the {id} path value; writes 404 and returns false on a bad id.
func adminID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.NotFound(w, r)
		return 0, false
	}
	return id, true
}

// validStatus returns status if it's in allowed, else "" (= all / unfiltered).
func validStatus(status string, allowed map[string]bool) string {
	if allowed[status] {
		return status
	}
	return ""
}

func countOrZero(st *store.Store, ctx context.Context, table, status string) int {
	n, err := st.CountByStatus(ctx, table, status)
	if err != nil {
		slog.Error("admin count error", "table", table, "err", err)
		return 0
	}
	return n
}

// userEmail returns the signed-in user's email for the admin chrome.
func userEmail(r *http.Request) string {
	if u := session.FromContext(r.Context()); u != nil {
		return u.Email
	}
	return ""
}

func isHTMX(r *http.Request) bool { return r.Header.Get("HX-Request") == "true" }

func render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	if err := c.Render(r.Context(), w); err != nil {
		slog.Error("render error", "err", err)
	}
}

func adminServerError(w http.ResponseWriter, r *http.Request, op string, err error) {
	slog.Error("admin error", "op", op, "err", err)
	http.Error(w, "Something went wrong.", http.StatusInternalServerError)
}

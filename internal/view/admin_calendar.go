package view

import (
	"time"

	"github.com/firefly-software-mt/advanced-template/internal/store"
)

// CalDay is one cell in the month grid.
type CalDay struct {
	ISO      string // YYYY-MM-DD
	Day      int    // day of month
	InMonth  bool   // belongs to the displayed month (vs leading/trailing spill)
	IsToday  bool
	Bookings []store.Booking // scheduled for this day
}

// CalMonth is the full calendar view model for one displayed month.
type CalMonth struct {
	Month       string // YYYY-MM displayed
	Label       string // "June 2026"
	Prev        string // YYYY-MM
	Next        string // YYYY-MM
	TodayMonth  string // current actual month, for the "Today" button
	IsCurrent   bool   // displayed month is the current month
	Weeks       [][]CalDay
	Unscheduled []store.Booking
}

// WeekdayLabels are the column headers, Sunday-first (US convention).
var WeekdayLabels = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

// gridDays is the number of cells in the month grid (6 weeks, always — keeps the
// grid height stable across months).
const gridDays = 42

// parseMonth parses a "2006-01" string to the first of that month (UTC, which is
// fine for date-only math). Falls back to now's month on a bad/empty value.
func parseMonth(month string, now time.Time) time.Time {
	if t, err := time.Parse("2006-01", month); err == nil {
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// gridStart is the Sunday on/just before the first of the month.
func gridStart(first time.Time) time.Time {
	return first.AddDate(0, 0, -int(first.Weekday()))
}

// CalendarQueryRange returns the [start, end) ISO date window the grid covers, so
// the handler can fetch exactly the scheduled bookings the grid will display.
func CalendarQueryRange(month string, now time.Time) (start, end string) {
	gs := gridStart(parseMonth(month, now))
	return gs.Format("2006-01-02"), gs.AddDate(0, 0, gridDays).Format("2006-01-02")
}

// BuildCalendar assembles the month grid from the scheduled + unscheduled
// bookings the handler fetched. now drives the "today" highlight.
func BuildCalendar(month string, scheduled, unscheduled []store.Booking, now time.Time) CalMonth {
	first := parseMonth(month, now)
	displayed := first.Format("2006-01")
	todayISO := now.Format("2006-01-02")
	todayMonth := now.Format("2006-01")

	byDate := make(map[string][]store.Booking, len(scheduled))
	for _, b := range scheduled {
		byDate[b.ScheduledDate] = append(byDate[b.ScheduledDate], b)
	}

	start := gridStart(first)
	weeks := make([][]CalDay, 0, gridDays/7)
	week := make([]CalDay, 0, 7)
	for i := 0; i < gridDays; i++ {
		d := start.AddDate(0, 0, i)
		iso := d.Format("2006-01-02")
		week = append(week, CalDay{
			ISO:      iso,
			Day:      d.Day(),
			InMonth:  d.Month() == first.Month() && d.Year() == first.Year(),
			IsToday:  iso == todayISO,
			Bookings: byDate[iso],
		})
		if len(week) == 7 {
			weeks = append(weeks, week)
			week = make([]CalDay, 0, 7)
		}
	}

	return CalMonth{
		Month:       displayed,
		Label:       first.Format("January 2006"),
		Prev:        first.AddDate(0, -1, 0).Format("2006-01"),
		Next:        first.AddDate(0, 1, 0).Format("2006-01"),
		TodayMonth:  todayMonth,
		IsCurrent:   displayed == todayMonth,
		Weeks:       weeks,
		Unscheduled: unscheduled,
	}
}

// calHasScheduled reports whether any in-month day carries a booking — drives
// the phone agenda's empty state.
func calHasScheduled(cal CalMonth) bool {
	for _, week := range cal.Weeks {
		for _, d := range week {
			if d.InMonth && len(d.Bookings) > 0 {
				return true
			}
		}
	}
	return false
}

// agendaDayLabel renders a phone agenda day heading like "Mon, Jun 8".
func agendaDayLabel(d CalDay) string {
	if t, err := time.Parse("2006-01-02", d.ISO); err == nil {
		return t.Format("Mon, Jan 2")
	}
	return d.ISO
}

// agendaMeta is the secondary line on an agenda row: plan, plus the pickup
// window when one is set.
func agendaMeta(b store.Booking) string {
	m := planShort(b.Plan)
	if b.PickupWindow != "" {
		m += " · " + b.PickupWindow
	}
	return m
}

// bookingDotClass returns the full class string for a status dot — a small
// circle paired on two channels so the seven statuses stay distinguishable
// without relying on hue alone: HUE groups the lifecycle (ochre = needs
// scheduling, teal = booked & waiting, rust = in Chanté's hands, brown = done,
// faint = off), and FILL splits each hue's two steps (a ring is the earlier /
// in-flight step, a solid is the later / action-needed step). Used by both the
// grid chips and the legend, so a dot always means the same thing.
func bookingDotClass(status string) string {
	const base = "h-2.5 w-2.5 shrink-0 rounded-full border-[1.5px] "
	switch status {
	case "new":
		return base + "border-warning bg-warning"
	case "scheduled":
		return base + "border-success bg-transparent"
	case "picked_up":
		return base + "border-success bg-success"
	case "in_progress":
		return base + "border-accent bg-transparent"
	case "ready":
		return base + "border-accent bg-accent"
	case "delivered":
		return base + "border-foreground bg-foreground"
	default: // canceled / unknown
		return base + "border-muted-foreground bg-transparent"
	}
}

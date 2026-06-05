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

// bookingDotClass is the status dot color for a compact calendar chip.
func bookingDotClass(status string) string {
	switch status {
	case "new":
		return "bg-warning"
	case "scheduled", "picked_up", "in_progress":
		return "bg-primary"
	case "ready":
		return "bg-accent"
	case "delivered":
		return "bg-success"
	default: // canceled / unknown
		return "bg-muted-foreground"
	}
}

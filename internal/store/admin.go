package store

import (
	"context"
	"fmt"
)

// Admin read/update methods backing the dashboard. The public site only ever
// inserts these rows; the admin works them (status, schedule, notes).

// GetBooking returns a single booking by id.
func (s *Store) GetBooking(ctx context.Context, id int64) (Booking, error) {
	var b Booking
	err := scanBooking(s.db.QueryRowContext(ctx, "SELECT "+bookingColumns+" FROM bookings WHERE id = ?", id), &b)
	return b, err
}

// UpdateBooking sets the operator-managed fields on a booking.
func (s *Store) UpdateBooking(ctx context.Context, id int64, status, scheduledDate, pickupWindow, adminNotes string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE bookings SET status = ?, scheduled_date = ?, pickup_window = ?, admin_notes = ? WHERE id = ?",
		status, scheduledDate, pickupWindow, adminNotes, id,
	)
	return err
}

// GetQuestion returns a single question by id.
func (s *Store) GetQuestion(ctx context.Context, id int64) (Question, error) {
	var q Question
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, contact, topic, message, status, admin_notes, created_at FROM questions WHERE id = ?", id,
	).Scan(&q.ID, &q.Name, &q.Contact, &q.Topic, &q.Message, &q.Status, &q.AdminNotes, &q.CreatedAt)
	return q, err
}

// UpdateQuestionStatus sets the triage status + notes on a question.
func (s *Store) UpdateQuestionStatus(ctx context.Context, id int64, status, adminNotes string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE questions SET status = ?, admin_notes = ? WHERE id = ?", status, adminNotes, id,
	)
	return err
}

// GetWaitlistEntry returns a single waitlist entry by id.
func (s *Store) GetWaitlistEntry(ctx context.Context, id int64) (WaitlistEntry, error) {
	var e WaitlistEntry
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, location, contact, kind, notes, status, admin_notes, created_at FROM waitlist WHERE id = ?", id,
	).Scan(&e.ID, &e.Name, &e.Location, &e.Contact, &e.Kind, &e.Notes, &e.Status, &e.AdminNotes, &e.CreatedAt)
	return e, err
}

// UpdateWaitlistStatus sets the outreach status + notes on a waitlist entry.
func (s *Store) UpdateWaitlistStatus(ctx context.Context, id int64, status, adminNotes string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE waitlist SET status = ?, admin_notes = ? WHERE id = ?", status, adminNotes, id,
	)
	return err
}

// countableTables whitelists the tables CountByStatus may query (the table name
// can't be a bound parameter, so it must never come from user input).
var countableTables = map[string]bool{"bookings": true, "questions": true, "waitlist": true}

// CountByStatus counts rows in one of the managed tables, optionally filtered by
// status ("" counts all). table must be one of bookings|questions|waitlist.
func (s *Store) CountByStatus(ctx context.Context, table, status string) (int, error) {
	if !countableTables[table] {
		return 0, fmt.Errorf("store: refusing to count unknown table %q", table)
	}
	var n int
	err := s.db.QueryRowContext(ctx,
		fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE (? = '' OR status = ?)", table),
		status, status,
	).Scan(&n)
	return n, err
}

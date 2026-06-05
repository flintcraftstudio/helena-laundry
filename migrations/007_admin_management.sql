-- +goose Up
-- Management state for the admin dashboard. `status` is free text so the booking
-- vocabulary can grow (e.g. the calendar's lifecycle) without another migration.

-- Bookings become workable pickups: a real scheduled date, a window, a status
-- lifecycle (new|scheduled|picked_up|in_progress|ready|delivered|canceled), notes.
ALTER TABLE bookings ADD COLUMN status TEXT NOT NULL DEFAULT 'new';
ALTER TABLE bookings ADD COLUMN scheduled_date TEXT NOT NULL DEFAULT '';   -- YYYY-MM-DD
ALTER TABLE bookings ADD COLUMN pickup_window TEXT NOT NULL DEFAULT '';
ALTER TABLE bookings ADD COLUMN admin_notes TEXT NOT NULL DEFAULT '';
CREATE INDEX idx_bookings_status_date ON bookings (status, scheduled_date);

-- Inquiries: triage state (new|handled|archived) + private notes.
ALTER TABLE questions ADD COLUMN status TEXT NOT NULL DEFAULT 'new';
ALTER TABLE questions ADD COLUMN admin_notes TEXT NOT NULL DEFAULT '';

-- Waitlist: outreach state (new|contacted|archived) + private notes.
ALTER TABLE waitlist ADD COLUMN status TEXT NOT NULL DEFAULT 'new';
ALTER TABLE waitlist ADD COLUMN admin_notes TEXT NOT NULL DEFAULT '';

-- +goose Down
DROP INDEX idx_bookings_status_date;
ALTER TABLE bookings DROP COLUMN status;
ALTER TABLE bookings DROP COLUMN scheduled_date;
ALTER TABLE bookings DROP COLUMN pickup_window;
ALTER TABLE bookings DROP COLUMN admin_notes;
ALTER TABLE questions DROP COLUMN status;
ALTER TABLE questions DROP COLUMN admin_notes;
ALTER TABLE waitlist DROP COLUMN status;
ALTER TABLE waitlist DROP COLUMN admin_notes;

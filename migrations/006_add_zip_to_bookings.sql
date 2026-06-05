-- +goose Up
-- Pickup ZIP, captured separately so the service area is checkable on the
-- backend and demand outside it is sortable later (mirrors waitlist.location).
ALTER TABLE bookings ADD COLUMN zip TEXT NOT NULL DEFAULT '';
CREATE INDEX idx_bookings_zip ON bookings (zip);

-- +goose Down
DROP INDEX idx_bookings_zip;
ALTER TABLE bookings DROP COLUMN zip;

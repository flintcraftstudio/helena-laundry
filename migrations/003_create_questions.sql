-- +goose Up
-- Questions / feedback submissions from the contact form.
CREATE TABLE questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    contact TEXT NOT NULL,        -- phone or email, as the visitor typed it
    topic TEXT NOT NULL,          -- question | feedback | schedule
    message TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE questions;

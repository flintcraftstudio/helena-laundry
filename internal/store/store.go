package store

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Store wraps the database connection and provides query methods.
type Store struct {
	db *sql.DB
}

// New creates a new Store.
func New(db *sql.DB) *Store {
	return &Store{db: db}
}

// CreateSession inserts a new session row.
func (s *Store) CreateSession(ctx context.Context, token string, userID int64, expiresAt time.Time) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
		token, userID, expiresAt,
	)
	return err
}

// GetSession retrieves a valid (non-expired) session by token.
func (s *Store) GetSession(ctx context.Context, token string) (int64, time.Time, error) {
	var userID int64
	var expiresAt time.Time
	err := s.db.QueryRowContext(ctx,
		"SELECT user_id, expires_at FROM sessions WHERE id = ? AND expires_at > CURRENT_TIMESTAMP",
		token,
	).Scan(&userID, &expiresAt)
	return userID, expiresAt, err
}

// DeleteSession removes a session by token.
func (s *Store) DeleteSession(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE id = ?", token)
	return err
}

// GetUserByID retrieves a user's id and email by their ID.
func (s *Store) GetUserByID(ctx context.Context, id int64) (int64, string, error) {
	var userID int64
	var email string
	err := s.db.QueryRowContext(ctx,
		"SELECT id, email FROM users WHERE id = ?",
		id,
	).Scan(&userID, &email)
	return userID, email, err
}

// GetUserByEmail retrieves a user by email, returning id, email, and password hash.
func (s *Store) GetUserByEmail(ctx context.Context, email string) (int64, string, string, error) {
	var id int64
	var userEmail, passwordHash string
	err := s.db.QueryRowContext(ctx,
		"SELECT id, email, password_hash FROM users WHERE email = ?",
		email,
	).Scan(&id, &userEmail, &passwordHash)
	return id, userEmail, passwordHash, err
}

// CreateUser inserts a new user with a bcrypt-hashed password.
func (s *Store) CreateUser(ctx context.Context, email, password string) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	result, err := s.db.ExecContext(ctx,
		"INSERT INTO users (email, password_hash) VALUES (?, ?)",
		email, string(hash),
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// DeleteExpiredSessions removes all expired sessions.
func (s *Store) DeleteExpiredSessions(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE expires_at <= CURRENT_TIMESTAMP")
	return err
}

// Question is a question/feedback submission from the contact form.
type Question struct {
	ID        int64
	Name      string
	Contact   string
	Topic     string
	Message   string
	CreatedAt time.Time
}

// CreateQuestion stores a question/feedback submission.
func (s *Store) CreateQuestion(ctx context.Context, name, contact, topic, message string) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO questions (name, contact, topic, message) VALUES (?, ?, ?, ?)",
		name, contact, topic, message,
	)
	return err
}

// ListQuestions returns the most recent question submissions, newest first.
func (s *Store) ListQuestions(ctx context.Context, limit int) ([]Question, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, name, contact, topic, message, created_at FROM questions ORDER BY created_at DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Question
	for rows.Next() {
		var q Question
		if err := rows.Scan(&q.ID, &q.Name, &q.Contact, &q.Topic, &q.Message, &q.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, q)
	}
	return out, rows.Err()
}

// WaitlistEntry is an out-of-area waitlist signup.
type WaitlistEntry struct {
	ID        int64
	Name      string
	Location  string
	Contact   string
	Kind      string
	Notes     string
	CreatedAt time.Time
}

// CreateWaitlistEntry stores an out-of-area waitlist signup.
func (s *Store) CreateWaitlistEntry(ctx context.Context, name, location, contact, kind, notes string) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO waitlist (name, location, contact, kind, notes) VALUES (?, ?, ?, ?, ?)",
		name, location, contact, kind, notes,
	)
	return err
}

// ListWaitlist returns waitlist signups, newest first.
func (s *Store) ListWaitlist(ctx context.Context, limit int) ([]WaitlistEntry, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, name, location, contact, kind, notes, created_at FROM waitlist ORDER BY created_at DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []WaitlistEntry
	for rows.Next() {
		var e WaitlistEntry
		if err := rows.Scan(&e.ID, &e.Name, &e.Location, &e.Contact, &e.Kind, &e.Notes, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

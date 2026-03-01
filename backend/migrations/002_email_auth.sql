ALTER TABLE users
  ADD COLUMN password_hash TEXT,
  ADD COLUMN provider TEXT NOT NULL DEFAULT 'google';

-- ensure emails are unique for traditional auth and google sync
CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique ON users(email);

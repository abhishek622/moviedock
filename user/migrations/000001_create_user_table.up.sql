CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- use pgcrypto gen_random_uuid() if available
  email TEXT NOT NULL UNIQUE,
  encrypted_password TEXT NOT NULL,           -- bcrypt hash (never plain password)
  full_name TEXT,
  role TEXT NOT NULL DEFAULT 'user',     -- simple single role string; consider roles table if complex
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  timezone TEXT,                         -- optional user timezone preference (e.g. "Asia/Kolkata")
  metadata JSONB DEFAULT '{}'::jsonb,    -- flexible per-user metadata
  last_login TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
);

-- indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

-- trigger for updated_at
CREATE TRIGGER update_user_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
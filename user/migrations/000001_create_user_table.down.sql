DROP INDEX IF EXISTS idx_users_email;
DROP TRIGGER IF EXISTS update_user_updated_at ON users;
DROP TABLE IF EXISTS users;
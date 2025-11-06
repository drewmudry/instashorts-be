-- Drop sessions table and indexes
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP TABLE IF EXISTS sessions;

-- Drop oauth_accounts table and indexes
DROP INDEX IF EXISTS idx_oauth_accounts_provider_provider_id;
DROP INDEX IF EXISTS idx_oauth_accounts_user_id;
DROP TABLE IF EXISTS oauth_accounts;

-- Drop users table and indexes
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;


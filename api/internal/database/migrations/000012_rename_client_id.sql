-- +migrate Up
ALTER TABLE sandboxes RENAME COLUMN guest_session_id TO client_id;
ALTER TABLE audit_logs RENAME COLUMN client_token TO client_id;
DROP INDEX IF EXISTS idx_audit_logs_client_token;
CREATE INDEX IF NOT EXISTS idx_audit_logs_client_id ON audit_logs(client_id);

-- +migrate Down
ALTER TABLE sandboxes RENAME COLUMN client_id TO guest_session_id;
ALTER TABLE audit_logs RENAME COLUMN client_id TO client_token;
DROP INDEX IF EXISTS idx_audit_logs_client_id;
CREATE INDEX IF NOT EXISTS idx_audit_logs_client_token ON audit_logs(client_token);

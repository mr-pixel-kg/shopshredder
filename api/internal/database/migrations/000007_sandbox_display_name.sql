-- +goose Up
ALTER TABLE sandboxes ADD COLUMN IF NOT EXISTS display_name VARCHAR(255) DEFAULT '';

-- +goose Down
ALTER TABLE sandboxes DROP COLUMN IF EXISTS display_name;

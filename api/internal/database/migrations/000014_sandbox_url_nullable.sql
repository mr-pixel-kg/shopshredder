-- +goose Up
ALTER TABLE sandboxes ALTER COLUMN url DROP NOT NULL;

-- +goose Down
ALTER TABLE sandboxes ALTER COLUMN url SET NOT NULL;

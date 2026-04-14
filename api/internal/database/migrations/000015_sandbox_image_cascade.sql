-- +goose Up
ALTER TABLE sandboxes DROP CONSTRAINT IF EXISTS sandboxes_image_id_fkey;
ALTER TABLE sandboxes
    ADD CONSTRAINT sandboxes_image_id_fkey
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE sandboxes DROP CONSTRAINT IF EXISTS sandboxes_image_id_fkey;
ALTER TABLE sandboxes
    ADD CONSTRAINT sandboxes_image_id_fkey
    FOREIGN KEY (image_id) REFERENCES images(id);

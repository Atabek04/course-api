-- +goose UP

ALTER TABLE module_info
    ADD CONSTRAINT updated_at_after_created_at CHECK (updated_at >= created_at),
    ADD CONSTRAINT module_duration_check CHECK (module_duration > 5 AND module_duration <= 15);

-- +goose Down
ALTER TABLE module_info
    DROP CONSTRAINT IF EXISTS updated_at_after_created_at,
    DROP CONSTRAINT IF EXISTS module_duration_check;
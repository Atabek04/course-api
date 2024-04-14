ALTER TABLE module_info
    DROP CONSTRAINT IF EXISTS updated_at_after_created_at,
    DROP CONSTRAINT IF EXISTS module_duration_check;

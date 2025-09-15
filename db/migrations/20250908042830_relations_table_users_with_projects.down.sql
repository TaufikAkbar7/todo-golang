DO $$
BEGIN
    -- drop column owner_id
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'projects' AND column_name='owner_id'
    ) THEN
        ALTER TABLE projects DROP owner_id;
    END IF;
END $$;
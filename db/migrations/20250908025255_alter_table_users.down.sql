DO $$
BEGIN
    -- rename column name to username
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='username')
    AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='name') THEN
        ALTER TABLE users RENAME COLUMN username TO name;
    END IF;

    -- drop unique constraint
    IF EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_users_username' AND table_name = 'users'
    ) THEN
        ALTER TABLE users DROP CONSTRAINT uq_users_username;
    END IF;
END $$;

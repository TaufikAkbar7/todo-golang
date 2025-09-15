DO $$
BEGIN
    -- drop fk constraint
    IF EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_project_members_role_id' AND table_name = 'project_members'
    ) THEN
        ALTER TABLE project_members DROP CONSTRAINT fk_project_members_role_id;
    END IF;

    -- rename column role_id to role and change data type VARCHAR(100) NOT NULL
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'project_members'
          AND column_name = 'role_id'
    ) AND NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'project_members'
          AND column_name = 'role'
    ) THEN
        ALTER TABLE project_members
            RENAME COLUMN role_id TO role;

        ALTER TABLE project_members
            ALTER COLUMN role SET DATA TYPE VARCHAR(100),
            ALTER COLUMN role SET NOT NULL;
    END IF;
END $$;

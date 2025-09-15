-- rename column
ALTER TABLE project_members RENAME COLUMN role TO role_id;

-- takeout not null constraint
ALTER TABLE project_members 
    ALTER COLUMN role_id DROP NOT NULL;

-- change type data
ALTER TABLE project_members
    ALTER COLUMN role_id TYPE UUID
    USING role_id::uuid;

-- add fk constraint
ALTER TABLE project_members ADD CONSTRAINT fk_project_members_role_id FOREIGN KEY(role_id) REFERENCES roles(id);

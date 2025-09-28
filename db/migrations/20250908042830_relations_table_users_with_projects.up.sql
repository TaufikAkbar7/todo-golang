ALTER TABLE projects
    ADD owner_id UUID,
    ADD CONSTRAINT fk_projects_owner
        FOREIGN KEY(owner_id)
        REFERENCES users(id) ON DELETE CASCADE;
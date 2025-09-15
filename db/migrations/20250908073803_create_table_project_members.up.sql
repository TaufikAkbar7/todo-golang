CREATE TABLE project_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID,
    user_id       UUID,
    role varchar(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_project_members_project_id
      FOREIGN KEY(project_id)
      REFERENCES projects(id)
      ON DELETE RESTRICT,
    CONSTRAINT fk_project_members_user_id
      FOREIGN KEY(user_id)
      REFERENCES users(id)
      ON DELETE RESTRICT
);
CREATE TABLE task_assigness (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  task_id UUID,
  project_member_id UUID,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_task_assigness_project_member_id
    FOREIGN KEY(project_member_id)
    REFERENCES project_members(id)
    ON DELETE CASCADE,
  CONSTRAINT fk_task_assigness_task_id
    FOREIGN KEY(task_id)
    REFERENCES tasks(id)
    ON DELETE CASCADE
);

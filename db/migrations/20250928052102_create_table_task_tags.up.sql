CREATE TABLE task_tags (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  task_id UUID,
  tag_id UUID,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_task_tags_task_id
    FOREIGN KEY(task_id)
    REFERENCES tasks(id)
    ON DELETE CASCADE,
  CONSTRAINT fk_task_tags_tag_id
    FOREIGN KEY(tag_id)
    REFERENCES tags(id)
    ON DELETE CASCADE
)
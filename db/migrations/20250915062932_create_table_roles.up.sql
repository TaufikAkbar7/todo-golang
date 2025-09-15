CREATE TABLE roles (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(50) UNIQUE
);

INSERT INTO roles (name) VALUES ('owner'), ('admin'), ('member'), ('viewer');
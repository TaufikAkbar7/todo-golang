CREATE TABLE roles (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) UNIQUE
);

INSERT INTO roles (name) VALUES ('owner'), ('admin'), ('member'), ('viewer');
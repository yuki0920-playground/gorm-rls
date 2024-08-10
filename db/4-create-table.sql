CREATE TABLE project_tags (
    project_id VARCHAR(255) NOT NULL REFERENCES projects(id),
    id VARCHAR(255) PRIMARY KEY
);

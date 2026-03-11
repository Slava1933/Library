CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    discipline VARCHAR(200) NOT NULL,
    file_path VARCHAR(400) NOT NULL,
);
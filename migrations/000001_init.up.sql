CREATE TABLE IF NOT EXISTS disciplines (
    id SERIAL PRIMARY KEY,
    title VARCHAR(300) NOT NULL,
    UNIQUE(title)
);

CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    discipline_id INTEGER REFERENCES disciplines(id), 
    title VARCHAR(200) NOT NULL,
    file_path VARCHAR(400) NOT NULL,
);
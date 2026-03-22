ALTER TABLE documents DROP CONSTRAINT
documents_discipline_id_fkey;

ALTER TABLE documents
ADD CONSTRAINT
documents_discipline_id_fkey
FOREIGN KEY (discipline_id)
REFERENCES disciplines(id);
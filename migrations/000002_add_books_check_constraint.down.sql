ALTER TABLE books DROP CONSTRAINT IF EXISTS books_pages_check;

ALTER TABLE books ADD CONSTRAINT IF EXISTS books_year_check;

ALTER TABLE books ADD CONSTRAINT IF EXISTS books_genres_length_check;
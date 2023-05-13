-- Add up migration script here
ALTER TABLE ratings
ADD CONSTRAINT fk_ratings_movies
FOREIGN KEY (tconst) REFERENCES movies (tconst)
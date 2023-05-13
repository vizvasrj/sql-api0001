-- Add down migration script here
ALTER TABLE ratings
DROP CONSTRAINT fk_ratings_movies;

ALTER TABLE ratings
ADD CONSTRAINT fk_ratings_movies
FOREIGN KEY (tconst) REFERENCES movies (tconst)
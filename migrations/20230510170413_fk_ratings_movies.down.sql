-- Add down migration script here
ALTER TABLE ratings
DROP CONSTRAINT fk_ratings_movies;
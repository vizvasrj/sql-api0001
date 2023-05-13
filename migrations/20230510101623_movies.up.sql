-- Add up migration script here
CREATE TABLE movies (
    tconst VARCHAR(10) PRIMARY KEY,
    titleType VARCHAR(10),
    primaryTitle VARCHAR(255),
    runtimeMinutes INTEGER,
    genres VARCHAR(255)
)
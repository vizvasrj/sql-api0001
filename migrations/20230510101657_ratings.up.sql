-- Add up migration script here
CREATE TABLE ratings (
    tconst VARCHAR(10) PRIMARY KEY,
    averageRating DECIMAL(3,1),
    numVotes INTEGER
)
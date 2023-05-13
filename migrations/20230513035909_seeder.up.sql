-- Add up migration script here
CREATE TABLE seeder(
    id SERIAL PRIMARY KEY,
    file_name VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
-- 000001_create_stamps_table.up.sql
CREATE TABLE IF NOT EXISTS stamps (
    id SERIAL PRIMARY KEY,
    blockheight INTEGER NOT NULL,
    stamp TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
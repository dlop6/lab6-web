CREATE TABLE IF NOT EXISTS matches (
    ID SERIAL PRIMARY KEY,
    Team1 VARCHAR(100) NOT NULL,
    Team2 VARCHAR(100) NOT NULL,
    Score1 INTEGER DEFAULT 0,
    Score2 INTEGER DEFAULT 0,
    Date DATE NOT NULL
);
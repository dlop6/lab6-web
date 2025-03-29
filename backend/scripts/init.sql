CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    team1 VARCHAR(100) NOT NULL,
    team2 VARCHAR(100) NOT NULL,
    score1 INTEGER DEFAULT 0 CHECK (score1 >= 0),
    score2 INTEGER DEFAULT 0 CHECK (score2 >= 0),
    date DATE NOT NULL,
    yellow_cards INTEGER DEFAULT 0 CHECK (yellow_cards >= 0),
    red_cards INTEGER DEFAULT 0 CHECK (red_cards >= 0),
    extra_time INTEGER DEFAULT 0 CHECK (extra_time BETWEEN 0 AND 30),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_matches_date ON matches(date);
CREATE INDEX IF NOT EXISTS idx_matches_teams ON matches(team1, team2);
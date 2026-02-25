CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    token TEXT UNIQUE NOT NULL,
    user_id TEXT NOT NULL,
    rate_limit INTEGER NOT NULL,
    tier TEXT NOT NULL
);

INSERT INTO users (token, user_id, rate_limit, tier) VALUES
('gold-token', 'user_pro_1', 100, 'gold'),
('free-token', 'user_norm_2', 5, 'free');

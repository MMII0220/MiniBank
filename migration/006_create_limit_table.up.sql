CREATE TABLE IF NOT EXISTS limits (
    ID SERIAL PRIMARY KEY,
    user_id INT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    daily_amount FLOAT NULL DEFAULT 10000, -- TJS
    last_reset TIMESTAMPTZ DEFAULT NOW()
);

-- -- Добавляем стандартные лимиты для всех существующих пользователей
-- INSERT INTO limits (user_id, daily_amount, last_reset)
-- SELECT id, 1000.0, NOW() 
-- FROM users
-- ON CONFLICT (user_id) DO NOTHING;

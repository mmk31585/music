-- +goose Up
-- +goose StatementBegin

-- Seed daily challenges that auto-track streaming and likes.
-- valid_from/valid_until use CURRENT_DATE so they are always current when the migration runs.

INSERT INTO daily_challenges (title, description, challenge_type, target_count, xp_reward, valid_from, valid_until)
VALUES
    ('Daily Streamer',     'Listen to 10 tracks today',   'streams', 10,  50,  CURRENT_DATE, CURRENT_DATE + INTERVAL '1 day'),
    ('Music Explorer',     'Listen to 30 tracks today',   'streams', 30,  100, CURRENT_DATE, CURRENT_DATE + INTERVAL '1 day'),
    ('Weekly Listener',    'Listen to 100 tracks this week', 'streams', 100, 300, CURRENT_DATE, CURRENT_DATE + INTERVAL '7 days'),
    ('Fan Favorite',       'Like 5 tracks',               'likes',   5,   50,  CURRENT_DATE, CURRENT_DATE + INTERVAL '1 day'),
    ('Super Fan',          'Like 15 tracks',              'likes',   15,  150, CURRENT_DATE, CURRENT_DATE + INTERVAL '7 days');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM daily_challenges WHERE title IN (
    'Daily Streamer',
    'Music Explorer',
    'Weekly Listener',
    'Fan Favorite',
    'Super Fan'
);

-- +goose StatementEnd

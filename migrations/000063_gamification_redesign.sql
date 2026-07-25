-- +goose Up
-- +goose StatementBegin

ALTER TABLE level_definitions
    ADD COLUMN IF NOT EXISTS can_upload BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS can_auto_publish BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS upload_slots INT NOT NULL DEFAULT 0;

COMMENT ON COLUMN level_definitions.can_upload IS 'Whether users at this level can upload tracks';
COMMENT ON COLUMN level_definitions.can_auto_publish IS 'Whether uploads skip review queue';
COMMENT ON COLUMN level_definitions.upload_slots IS 'Number of concurrent upload slots at this level';

UPDATE level_definitions SET can_upload = FALSE, can_auto_publish = FALSE, upload_slots = 0;
UPDATE level_definitions SET can_upload = TRUE, upload_slots = 1 WHERE level >= 5;
UPDATE level_definitions SET upload_slots = 2 WHERE level >= 7;
UPDATE level_definitions SET upload_slots = 5 WHERE level >= 8;
UPDATE level_definitions SET can_auto_publish = TRUE, upload_slots = 10 WHERE level >= 9;
UPDATE level_definitions SET upload_slots = 20 WHERE level >= 12;
UPDATE level_definitions SET upload_slots = 50 WHERE level >= 15;

INSERT INTO badges (name, description, icon_url, category, rarity, criteria, xp_reward, is_hidden) VALUES
    ('First Upload',        'Upload your first track',              '/badges/first-upload.svg',        'upload', 'common',   '{"type":"upload_count","count":1}',          100, FALSE),
    ('Prolific Uploader',   'Upload 5 tracks',                      '/badges/upload-5.svg',            'upload', 'uncommon', '{"type":"upload_count","count":5}',          250, FALSE),
    ('Upload Master',        'Upload 25 tracks',                     '/badges/upload-25.svg',           'upload', 'rare',     '{"type":"upload_count","count":25}',         500, FALSE),
    ('Upload Legend',        'Upload 100 tracks',                    '/badges/upload-100.svg',          'upload', 'epic',     '{"type":"upload_count","count":100}',       1000, FALSE),
    ('First Contribution',  'Make your first metadata contribution', '/badges/first-contribution.svg',  'community', 'common', '{"type":"contribution_count","count":1}',  50, FALSE),
    ('Community Helper',     'Make 10 accepted contributions',       '/badges/contribution-10.svg',     'community', 'uncommon','{"type":"contribution_count","count":10}', 200, FALSE),
    ('Community Pillar',     'Make 50 accepted contributions',       '/badges/contribution-50.svg',     'community', 'rare',   '{"type":"contribution_count","count":50}', 500, FALSE),
    ('Trusted Reviewer',     'Review 25 uploads',                    '/badges/trusted-reviewer.svg',    'moderation', 'rare',  '{"type":"review_count","count":25}',       300, FALSE),
    ('Track Published',     'Your first upload was published',      '/badges/first-published.svg',     'upload', 'common',  '{"type":"track_published","count":1}',     150, FALSE),
    ('7-Day Streak',         'Contribute 7 days in a row',           '/badges/streak-7.svg',            'activity', 'uncommon','{"type":"streak_days","count":7}',        200, FALSE),
    ('30-Day Streak',        'Contribute 30 days in a row',          '/badges/streak-30.svg',           'activity', 'rare',   '{"type":"streak_days","count":30}',        500, FALSE),
    ('Rising Star',          'Reach level 5',                        '/badges/level-5.svg',             'milestone', 'uncommon','{"type":"level_reached","level":5}',     200, FALSE),
    ('Music Expert',         'Reach level 10',                       '/badges/level-10.svg',            'milestone', 'rare',   '{"type":"level_reached","level":10}',      500, FALSE),
    ('Music Master',         'Reach level 15',                       '/badges/level-15.svg',            'milestone', 'epic',   '{"type":"level_reached","level":15}',     1000, FALSE),
    ('Trusted Creator',      'Earn auto-publish privilege',          '/badges/auto-publish.svg',        'upload', 'rare',     '{"type":"auto_publish_earned"}',           300, FALSE)
;

INSERT INTO daily_challenges (title, description, challenge_type, target_count, xp_reward, badge_reward_id, is_active, valid_from, valid_until) VALUES
    ('Upload Today',      'Upload at least 1 track today',     'upload', 1, 50, NULL, TRUE, NOW(), NOW() + INTERVAL '1 day'),
    ('Help the Community', 'Make 3 metadata contributions today', 'contribute', 3, 75, NULL, TRUE, NOW(), NOW() + INTERVAL '1 day'),
    ('Review Queue',      'Review 5 pending uploads today',    'review', 5, 100, NULL, TRUE, NOW(), NOW() + INTERVAL '1 day');

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'xp_transactions_source_check') THEN
        ALTER TABLE xp_transactions DROP CONSTRAINT xp_transactions_source_check;
    END IF;
END $$;

ALTER TABLE xp_transactions
    ADD CONSTRAINT xp_transactions_source_check
    CHECK (source IN (
        'stream', 'like', 'share', 'contribution', 'translation', 'login',
        'challenge_stream', 'challenge_like', 'challenge_share',
        'upload', 'upload_published', 'contribution_accepted',
        'review', 'streak', 'badge', 'level_up',
        'co_upload', 'club_submit'
    ));

CREATE INDEX IF NOT EXISTS idx_xp_transactions_source ON xp_transactions(source);
CREATE INDEX IF NOT EXISTS idx_xp_transactions_created ON xp_transactions(created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_xp_transactions_created;
DROP INDEX IF EXISTS idx_xp_transactions_source;

ALTER TABLE xp_transactions DROP CONSTRAINT IF EXISTS xp_transactions_source_check;

ALTER TABLE xp_transactions
    ADD CONSTRAINT xp_transactions_source_check
    CHECK (source IN (
        'stream', 'like', 'share', 'contribution', 'translation', 'login',
        'challenge_stream', 'challenge_like', 'challenge_share'
    ));

DELETE FROM daily_challenges WHERE id IN ('upload_1_daily', 'contribute_3_daily', 'review_5_daily');
DELETE FROM badges WHERE id IN (
    'first_upload', 'upload_5', 'upload_25', 'upload_100',
    'first_contribution', 'contribution_10', 'contribution_50',
    'trusted_reviewer', 'first_track_published',
    'streak_7', 'streak_30', 'level_5', 'level_10', 'level_15', 'auto_publish'
);

-- +goose StatementEnd
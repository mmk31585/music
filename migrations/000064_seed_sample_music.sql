-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$
DECLARE
  artist_id UUID := gen_random_uuid();
  track_id UUID := gen_random_uuid();
  existing_artist_id UUID;
  existing_track_id UUID;
BEGIN
  -- Seed a sample artist (placeholder, no real person). Idempotent on slug.
  INSERT INTO artists (id, name, slug, bio, is_verified, monthly_listeners)
  VALUES (gen_random_uuid(), 'Sample Artist', 'sample-artist',
          'Placeholder artist used for end-to-end smoke testing of playback and lyrics.',
          FALSE, 0)
  ON CONFLICT (slug) DO UPDATE
  SET name = EXCLUDED.name,
      bio = EXCLUDED.bio,
      is_verified = EXCLUDED.is_verified,
      monthly_listeners = EXCLUDED.monthly_listeners;

  -- Resolve the (possibly pre-existing) artist id deterministically.
  SELECT id INTO existing_artist_id FROM artists WHERE slug = 'sample-artist';

  -- Seed a sample track (placeholder audio URL; no real media file).
  -- Idempotent on (artist_id, slug); updates fields on re-apply.
  INSERT INTO tracks (id, artist_id, title, slug, duration_seconds, audio_url, cover_url, is_public)
  VALUES (gen_random_uuid(), existing_artist_id, 'Sample Track', 'sample-track', 240,
          'https://example.com/placeholder/sample-track.mp3',
          'https://example.com/placeholder/sample-track-cover.jpg',
          TRUE)
  ON CONFLICT (artist_id, slug) DO UPDATE
  SET title = EXCLUDED.title,
      duration_seconds = EXCLUDED.duration_seconds,
      audio_url = EXCLUDED.audio_url,
      cover_url = EXCLUDED.cover_url,
      is_public = EXCLUDED.is_public;

  -- Resolve the (possibly pre-existing) track id deterministically.
  SELECT id INTO existing_track_id FROM tracks WHERE slug = 'sample-track';

  -- Seed sample lyrics for the track (original, non-copyrighted placeholder text).
  -- Idempotent on (track_id, language); upserts on re-apply.
  INSERT INTO lyrics (track_id, language, type, content, source, confidence_score, created_at, updated_at)
  VALUES (existing_track_id, 'fa', 'plain',
          'صبح خوب است، روزی تازه است،' || E'\n' ||
          'موزیک در دل، Heart resonant.' || E'\n' ||
          'این affect تنها یک نمونه است.',
          'manual', NULL, NOW(), NOW())
  ON CONFLICT (track_id, language) DO UPDATE
  SET type = EXCLUDED.type,
      content = EXCLUDED.content,
      source = EXCLUDED.source,
      confidence_score = EXCLUDED.confidence_score,
      updated_at = NOW();
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM lyrics WHERE track_id IN (
  SELECT id FROM tracks WHERE slug = 'sample-track'
);
DELETE FROM tracks WHERE slug = 'sample-track';
DELETE FROM artists WHERE slug = 'sample-artist';

-- +goose StatementEnd

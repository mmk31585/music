-- +goose Up
-- +goose StatementBegin

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Pop') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Pop', 'pop');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Traditional') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Traditional', 'traditional');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Classical') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Classical', 'classical');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Rock') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Rock', 'rock');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Hip Hop') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Hip Hop', 'hip-hop');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Rap') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Rap', 'rap');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Electronic') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Electronic', 'electronic');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Folk') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Folk', 'folk');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Jazz') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Jazz', 'jazz');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Blues') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Blues', 'blues');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'R&B') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'R&B', 'rnb');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Soul') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Soul', 'soul');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Metal') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Metal', 'metal');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Fusion') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Fusion', 'fusion');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Reggae') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Reggae', 'reggae');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Latin') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Latin', 'latin');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Country') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Country', 'country');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Instrumental') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Instrumental', 'instrumental');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Lo-Fi') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Lo-Fi', 'lofi');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Ambient') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Ambient', 'ambient');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Pop Rock') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Pop Rock', 'pop-rock');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Indie') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Indie', 'indie');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Funk') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Funk', 'funk');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'Dance') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'Dance', 'dance');
END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'New Age') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'New Age', 'new-age');
END IF;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM track_genres WHERE genre_id IN (SELECT id FROM genres);
DELETE FROM genres;

-- +goose StatementEnd

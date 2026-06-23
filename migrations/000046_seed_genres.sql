-- +goose Up
-- +goose StatementBegin

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'پاپ') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'پاپ', 'pop');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'سنتی') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'سنتی', 'traditional');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'موسیقی کلاسیک') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'موسیقی کلاسیک', 'classical');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'راک') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'راک', 'rock');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'هیپ هاپ') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'هیپ هاپ', 'hip-hop');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'رپ') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'رپ', 'rap');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'الکترونیک') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'الکترونیک', 'electronic');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'محلی') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'محلی', 'folk');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'جاز') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'جاز', 'jazz');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'بلوز') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'بلوز', 'blues');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'R&B') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'R&B', 'rnb');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'سول') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'سول', 'soul');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'متال') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'متال', 'metal');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'تلفیقی') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'تلفیقی', 'fusion');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'رگی') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'رگی', 'reggae');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'لاتین') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'لاتین', 'latin');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'کانتری') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'کانتری', 'country');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'بی کلام') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'بی کلام', 'instrumental');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'لایت') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'لایت', 'lofi');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'آمبینت') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'آمبینت', 'ambient');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'پاپ راک') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'پاپ راک', 'pop-rock');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'ایندی') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'ایندی', 'indie');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'فانک') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'فانک', 'funk');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'دنس') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'دنس', 'dance');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM genres WHERE name = 'نیو ایج') THEN
        INSERT INTO genres (id, name, slug) VALUES (gen_random_uuid(), 'نیو ایج', 'new-age');
    END IF;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM track_genres WHERE genre_id IN (SELECT id FROM genres);
DELETE FROM genres;

-- +goose StatementEnd

package finalization

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/gif"
	_ "image/png"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/image/draw"

	"music/internal/modules/catalog/common"
	platformstorage "music/internal/platform/storage"
)

const (
	CatalogAudioDir = "catalog-audio"
	CatalogCoverDir = "catalog-covers"
)

type FinalizeResult struct {
	ArtistID      string `json:"artistId"`
	AlbumID       string `json:"albumId,omitempty"`
	TrackID       string `json:"trackId"`
	AudioURL      string `json:"audioUrl"`
	CoverURL      string `json:"coverUrl,omitempty"`
	CoverThumbURL string `json:"coverThumbUrl,omitempty"`
	CoverMedURL   string `json:"coverMedUrl,omitempty"`
}

type DraftMetadata struct {
	Artist ArtistMeta `json:"artist"`
	Album  AlbumMeta  `json:"album"`
	Track  TrackMeta  `json:"track"`
}

type ArtistMeta struct {
	Action          string `json:"action"`
	ExistingID      string `json:"existingId,omitempty"`
	Name            string `json:"name"`
	Bio             string `json:"bio,omitempty"`
	ImageURL        string `json:"imageUrl,omitempty"`
	Country         string `json:"country,omitempty"`
	MusicBrainzMBID string `json:"musicbrainzMbid,omitempty"`
}

type AlbumMeta struct {
	Action               string `json:"action"`
	ExistingID           string `json:"existingId,omitempty"`
	Title                string `json:"title"`
	ReleaseYear          int    `json:"releaseYear,omitempty"`
	Genre                string `json:"genre,omitempty"`
	CoverURL             string `json:"coverUrl,omitempty"`
	MusicBrainzReleaseID string `json:"musicbrainzReleaseId,omitempty"`
}

type TrackMeta struct {
	Title            string `json:"title"`
	TrackNumber      int    `json:"trackNumber,omitempty"`
	DurationSeconds  int    `json:"durationSeconds"`
	Genre            string `json:"genre,omitempty"`
	Lyrics           string `json:"lyrics,omitempty"`
	Explicit         bool   `json:"explicit"`
	SpotifyPreviewURL string `json:"spotifyPreviewUrl,omitempty"`
	CoverURL         string `json:"coverUrl,omitempty"`
}

type DraftRef struct {
	ID         string
	Status     string
	FinalMeta  string
	FilePath   string
	Format     string
}

type AssetRef struct {
	AssetType string
	URL       string
}

type Service struct {
	db      *sqlx.DB
	storage platformstorage.Storage
	logger  *zap.Logger
}

func NewService(db *sqlx.DB, storage platformstorage.Storage, logger *zap.Logger) *Service {
	return &Service{db: db, storage: storage, logger: logger}
}

func (s *Service) Finalize(ctx context.Context, draftID, draftStatus, draftFilepath, draftFormat, finalMetaJSON string) (*FinalizeResult, error) {
	if draftStatus != "accepted" {
		return nil, fmt.Errorf("draft must be in accepted status, got %s", draftStatus)
	}
	if finalMetaJSON == "" {
		return nil, fmt.Errorf("draft has no final metadata")
	}

	var final DraftMetadata
	if err := json.Unmarshal([]byte(finalMetaJSON), &final); err != nil {
		return nil, fmt.Errorf("unmarshal final metadata: %w", err)
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	artistID, err := s.upsertArtistTx(ctx, tx, &final.Artist)
	if err != nil {
		return nil, fmt.Errorf("artist: %w", err)
	}
	s.logger.Info("artist resolved", zap.String("artist_id", artistID.String()))

	var albumID uuid.UUID
	if final.Album.Action != "skip" {
		albumID, err = s.upsertAlbumTx(ctx, tx, &final.Album, artistID)
		if err != nil {
			return nil, fmt.Errorf("album: %w", err)
		}
		s.logger.Info("album resolved", zap.String("album_id", albumID.String()))
	}

	coverURL := final.Track.CoverURL
	if coverURL == "" {
		coverURL = final.Album.CoverURL
	}

	audioKey := draftFilepath
	lastSlash := strings.LastIndex(audioKey, "/")
	audioFilename := audioKey
	if lastSlash != -1 {
		audioFilename = audioKey[lastSlash+1:]
	}
	dstAudioKey := fmt.Sprintf("%s/%s/%s", CatalogAudioDir, artistID.String(), audioFilename)

	if err := s.storage.Copy(ctx, audioKey, dstAudioKey); err != nil {
		return nil, fmt.Errorf("copy audio: %w", err)
	}

	copiedKeys := []string{dstAudioKey}
	defer func() {
		if err != nil {
			for _, key := range copiedKeys {
				if delErr := s.storage.Delete(ctx, key); delErr != nil {
					s.logger.Warn("failed to clean up copied file", zap.String("key", key), zap.Error(delErr))
				}
			}
		}
	}()

	audioURL, err := s.storage.GetURL(ctx, dstAudioKey)
	if err != nil {
		return nil, fmt.Errorf("audio url: %w", err)
	}

	result := &FinalizeResult{
		ArtistID: artistID.String(),
		AlbumID:  albumID.String(),
	}

	assets, err := s.getAssetsByDraftID(ctx, draftID)
	if err != nil {
		return nil, fmt.Errorf("get assets: %w", err)
	}

	for _, asset := range assets {
		if asset.AssetType == "cover" {
			coverKey := extractStorageKey(asset.URL, s.getBaseURLPrefix(ctx))
			if idx := strings.LastIndex(coverKey, "."); idx != -1 {
				ext := coverKey[idx:]
				baseCoverKey := fmt.Sprintf("%s/%s/cover", CatalogCoverDir, albumID.String())

				dstCoverKey := baseCoverKey + ext
				if err := s.storage.Copy(ctx, coverKey, dstCoverKey); err != nil {
					s.logger.Warn("failed to copy cover art", zap.String("src", coverKey), zap.Error(err))
				} else {
					copiedKeys = append(copiedKeys, dstCoverKey)
					if cvURL, cvErr := s.storage.GetURL(ctx, dstCoverKey); cvErr == nil {
						coverURL = cvURL
					}
				}

				coverData, err := s.storage.Get(ctx, coverKey)
				if err == nil && len(coverData) > 0 {
					thumbURL, medURL, resizeErr := s.resizeAndStoreCover(ctx, coverData, ext, baseCoverKey)
					if resizeErr != nil {
						s.logger.Warn("cover resize failed", zap.Error(resizeErr))
					} else {
						result.CoverThumbURL = thumbURL
						result.CoverMedURL = medURL
						if thumbURL != "" {
							thumbKey := extractStorageKey(thumbURL, s.getBaseURLPrefix(ctx))
							copiedKeys = append(copiedKeys, thumbKey)
						}
						if medURL != "" {
							medKey := extractStorageKey(medURL, s.getBaseURLPrefix(ctx))
							copiedKeys = append(copiedKeys, medKey)
						}
					}
				}
			} else {
				coverKey = fmt.Sprintf("%s/%s/cover", CatalogCoverDir, albumID.String()) + ".jpg"
				if err := s.storage.Copy(ctx, coverKey, coverKey); err != nil {
					s.logger.Warn("failed to copy cover art", zap.String("src", coverKey), zap.Error(err))
				} else {
					copiedKeys = append(copiedKeys, coverKey)
					if cvURL, cvErr := s.storage.GetURL(ctx, coverKey); cvErr == nil {
						coverURL = cvURL
					}
				}
			}
			break
		}
	}

	explicit := final.Track.Explicit
	trackNumber := 1
	if final.Track.TrackNumber > 0 {
		trackNumber = final.Track.TrackNumber
	}

	genres := splitGenres(final.Track.Genre)
	genreIDs := make([]uuid.UUID, 0, len(genres))
	for _, g := range genres {
		gid, err := s.ensureGenreTx(ctx, tx, g)
		if err != nil {
			return nil, fmt.Errorf("genre %q: %w", g, err)
		}
		genreIDs = append(genreIDs, gid)
	}

	var albumIDPtr *uuid.UUID
	if albumID != uuid.Nil {
		albumIDPtr = &albumID
	}

	trackID, err := s.insertTrackTx(ctx, tx, insertTrackParams{
		Title:           final.Track.Title,
		ArtistID:        artistID,
		AlbumID:         albumIDPtr,
		DurationSeconds: final.Track.DurationSeconds,
		TrackNumber:     trackNumber,
		Explicit:        explicit,
		AudioURL:        audioURL,
		CoverURL:        coverURL,
		GenreIDs:        genreIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("track: %w", err)
	}
	s.logger.Info("track created", zap.String("track_id", trackID.String()))

	if final.Track.Lyrics != "" {
		lyricsType := "plain"
		if looksLikeLRC(final.Track.Lyrics) {
			lyricsType = "lrc"
		}
		if err := s.insertLyricsTx(ctx, tx, trackID, "en", lyricsType, final.Track.Lyrics); err != nil {
			return nil, fmt.Errorf("lyrics: %w", err)
		}
		s.logger.Info("lyrics inserted", zap.String("track_id", trackID.String()), zap.String("type", lyricsType))
	}

	if err := s.updateDraftResultTx(ctx, tx, draftID, artistID, albumID, trackID); err != nil {
		return nil, fmt.Errorf("update draft: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	result.AudioURL = audioURL
	result.CoverURL = coverURL
	result.TrackID = trackID.String()

	return result, nil
}

type insertTrackParams struct {
	Title           string
	ArtistID        uuid.UUID
	AlbumID         *uuid.UUID
	DurationSeconds int
	TrackNumber     int
	Explicit        bool
	AudioURL        string
	CoverURL        string
	GenreIDs        []uuid.UUID
}

func (s *Service) insertTrackTx(ctx context.Context, tx *sqlx.Tx, p insertTrackParams) (uuid.UUID, error) {
	slug := common.Slugify(p.Title)

	var trackID uuid.UUID
	err := tx.GetContext(ctx, &trackID, `
		INSERT INTO tracks (
			artist_id, album_id, title, slug,
			duration_seconds, track_number, explicit,
			audio_url, cover_url, is_public
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, TRUE)
		RETURNING id
	`,
		p.ArtistID,
		p.AlbumID,
		p.Title,
		slug,
		p.DurationSeconds,
		p.TrackNumber,
		p.Explicit,
		p.AudioURL,
		p.CoverURL,
	)
	if err != nil {
		return uuid.Nil, common.MapPGError(err)
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO track_artists (track_id, artist_id, role, position)
		VALUES ($1, $2, 'primary', 1)
	`, trackID, p.ArtistID); err != nil {
		return uuid.Nil, common.MapPGError(err)
	}

	for _, gid := range p.GenreIDs {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO track_genres (track_id, genre_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, trackID, gid); err != nil {
			return uuid.Nil, common.MapPGError(err)
		}
	}

	return trackID, nil
}

func (s *Service) upsertArtistTx(ctx context.Context, tx *sqlx.Tx, meta *ArtistMeta) (uuid.UUID, error) {
	if meta.Action == "link" && meta.ExistingID != "" {
		uid, err := common.ParseUUID(meta.ExistingID)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid existing artist ID: %w", err)
		}
		return uid, nil
	}

	slug := common.Slugify(meta.Name)
	var artistID uuid.UUID
	err := tx.GetContext(ctx, &artistID, `
		INSERT INTO artists (name, slug, bio, image_url, country)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`,
		meta.Name,
		slug,
		strPtr(meta.Bio),
		strPtr(meta.ImageURL),
		strPtr(meta.Country),
	)
	if err != nil {
		return uuid.Nil, common.MapPGError(err)
	}
	return artistID, nil
}

func (s *Service) upsertAlbumTx(ctx context.Context, tx *sqlx.Tx, meta *AlbumMeta, artistID uuid.UUID) (uuid.UUID, error) {
	if meta.Action == "link" && meta.ExistingID != "" {
		uid, err := common.ParseUUID(meta.ExistingID)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid existing album ID: %w", err)
		}
		return uid, nil
	}

	slug := common.Slugify(meta.Title)
	albumType := "album"

	var releaseDate interface{}
	if meta.ReleaseYear > 0 {
		rd := time.Date(meta.ReleaseYear, 1, 1, 0, 0, 0, 0, time.UTC)
		releaseDate = &rd
	}

	var albumID uuid.UUID
	err := tx.GetContext(ctx, &albumID, `
		INSERT INTO albums (artist_id, title, slug, cover_url, release_date, album_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (artist_id, slug) DO UPDATE SET title = EXCLUDED.title
		RETURNING id
	`,
		artistID,
		meta.Title,
		slug,
		strPtr(meta.CoverURL),
		releaseDate,
		albumType,
	)
	if err != nil {
		return uuid.Nil, common.MapPGError(err)
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO album_artists (album_id, artist_id, role, position)
		VALUES ($1, $2, 'primary', 1)
		ON CONFLICT DO NOTHING
	`, albumID, artistID); err != nil {
		return uuid.Nil, common.MapPGError(err)
	}

	return albumID, nil
}

func (s *Service) ensureGenreTx(ctx context.Context, tx *sqlx.Tx, name string) (uuid.UUID, error) {
	slug := common.Slugify(name)

	var id uuid.UUID
	err := tx.GetContext(ctx, &id, `SELECT id FROM genres WHERE slug = $1`, slug)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return uuid.Nil, err
	}

	err = tx.GetContext(ctx, &id, `
		INSERT INTO genres (name, slug)
		VALUES ($1, $2)
		RETURNING id
	`, name, slug)
	if err != nil {
		return uuid.Nil, common.MapPGError(err)
	}
	return id, nil
}

func (s *Service) updateDraftResultTx(ctx context.Context, tx *sqlx.Tx, draftID string, artistID, albumID, trackID uuid.UUID) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE ingestion_drafts
		SET status = $1, artist_id = $2, album_id = $3, track_id = $4, updated_at = NOW()
		WHERE id = $5
	`, "published", artistID, albumID, trackID, draftID)
	return err
}

type draftAsset struct {
	AssetType string `db:"asset_type"`
	URL       string `db:"url"`
}

func (s *Service) getAssetsByDraftID(ctx context.Context, draftID string) ([]draftAsset, error) {
	var assets []draftAsset
	err := s.db.SelectContext(ctx, &assets, `
		SELECT asset_type, url
		FROM ingestion_draft_assets
		WHERE draft_id = $1
		ORDER BY created_at ASC
	`, draftID)
	return assets, err
}

func (s *Service) getBaseURLPrefix(ctx context.Context) string {
	url, err := s.storage.GetURL(ctx, "")
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(url, "/")
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func splitGenres(genre string) []string {
	if genre == "" {
		return nil
	}
	parts := strings.Split(genre, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func extractStorageKey(url, baseURL string) string {
	key := strings.TrimPrefix(url, baseURL+"/")
	key = strings.TrimPrefix(key, "/")
	return key
}

const (
	coverThumbSize = 100
	coverMedSize   = 300
)

func (s *Service) resizeAndStoreCover(ctx context.Context, data []byte, ext, baseKey string) (thumbURL, medURL string, err error) {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", "", fmt.Errorf("decode cover: %w", err)
	}

	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w <= coverMedSize && h <= coverMedSize {
		return "", "", nil
	}

	if w > coverThumbSize || h > coverThumbSize {
		thumbKey := baseKey + "_thumb" + ext
		thumb := image.NewRGBA(image.Rect(0, 0, coverThumbSize, coverThumbSize))
		draw.ApproxBiLinear.Scale(thumb, thumb.Bounds(), src, bounds, draw.Src, nil)
		if uErr := s.uploadImage(ctx, thumbKey, thumb); uErr != nil {
			s.logger.Warn("failed to upload thumb cover", zap.Error(uErr))
		} else {
			if tURL, tErr := s.storage.GetURL(ctx, thumbKey); tErr == nil {
				thumbURL = tURL
			}
		}
	}

	if w > coverMedSize || h > coverMedSize {
		medKey := baseKey + "_med" + ext
		med := image.NewRGBA(image.Rect(0, 0, coverMedSize, coverMedSize))
		draw.ApproxBiLinear.Scale(med, med.Bounds(), src, bounds, draw.Src, nil)
		if uErr := s.uploadImage(ctx, medKey, med); uErr != nil {
			s.logger.Warn("failed to upload med cover", zap.Error(uErr))
		} else {
			if mURL, mErr := s.storage.GetURL(ctx, medKey); mErr == nil {
				medURL = mURL
			}
		}
	}

	return thumbURL, medURL, nil
}

func (s *Service) uploadImage(ctx context.Context, key string, img image.Image) error {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85}); err != nil {
		return err
	}
	return s.storage.Upload(ctx, key, &buf, int64(buf.Len()), "image/jpeg")
}

func (s *Service) insertLyricsTx(ctx context.Context, tx *sqlx.Tx, trackID uuid.UUID, language, lyricsType, content string) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO lyrics (track_id, language, type, content)
		VALUES ($1, $2, $3, $4)
	`, trackID, language, lyricsType, content)
	if err != nil {
		return err
	}
	return nil
}

func looksLikeLRC(content string) bool {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		return strings.HasPrefix(line, "[")
	}
	return false
}

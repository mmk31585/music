package social

import (
	"context"

	"github.com/google/uuid"
)

// --- Room Queue Voting ---

func (r *repository) SuggestTrack(ctx context.Context, c *QueueCandidate) error {
	c.ID = uuid.New()
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO room_queue_candidates (id, room_id, track_id, suggested_by, vote_count, created_at)
		VALUES (:id, :room_id, :track_id, :suggested_by, :vote_count, :created_at)
		ON CONFLICT (room_id, track_id) DO NOTHING
	`, c)
	return err
}

func (r *repository) GetCandidates(ctx context.Context, roomID uuid.UUID) ([]QueueCandidate, error) {
	items := make([]QueueCandidate, 0)
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, room_id, track_id, suggested_by, vote_count, created_at
		FROM room_queue_candidates
		WHERE room_id = $1
		ORDER BY vote_count DESC, created_at ASC
	`, roomID)
	return items, err
}

func (r *repository) GetCandidateByID(ctx context.Context, candidateID uuid.UUID) (*QueueCandidate, error) {
	var c QueueCandidate
	err := r.db.GetContext(ctx, &c, `
		SELECT id, room_id, track_id, suggested_by, vote_count, created_at
		FROM room_queue_candidates
		WHERE id = $1
	`, candidateID)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) CastVoteTx(ctx context.Context, candidateID, userID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.GetContext(ctx, &exists, `
		SELECT EXISTS(SELECT 1 FROM room_queue_votes WHERE candidate_id = $1 AND user_id = $2)
	`, candidateID, userID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO room_queue_votes (candidate_id, user_id)
		VALUES ($1, $2)
	`, candidateID, userID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE room_queue_candidates
		SET vote_count = vote_count + 1
		WHERE id = $1
	`, candidateID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) RemoveVoteTx(ctx context.Context, candidateID, userID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `
		DELETE FROM room_queue_votes
		WHERE candidate_id = $1 AND user_id = $2
	`, candidateID, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE room_queue_candidates
		SET vote_count = GREATEST(vote_count - 1, 0)
		WHERE id = $1
	`, candidateID)

	return tx.Commit()
}

func (r *repository) HasVoted(ctx context.Context, candidateID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS(SELECT 1 FROM room_queue_votes WHERE candidate_id = $1 AND user_id = $2)
	`, candidateID, userID)
	return exists, err
}

func (r *repository) GetNowPlaying(ctx context.Context, roomID uuid.UUID) (*RoomNowPlaying, error) {
	var np RoomNowPlaying
	err := r.db.GetContext(ctx, &np, `
		SELECT room_id, track_id, started_at, suggested_by, source
		FROM room_now_playing
		WHERE room_id = $1
	`, roomID)
	if err != nil {
		return nil, err
	}
	return &np, nil
}

func (r *repository) GetNowPlayingWithTrack(ctx context.Context, roomID uuid.UUID) (*NowPlayingResponse, error) {
	var row nowPlayingTrackRow
	err := r.db.GetContext(ctx, &row, `
		SELECT
			np.room_id, np.track_id, np.started_at, np.suggested_by, np.source,
			t.title, t.duration_seconds, t.cover_url,
			a.name AS artist_name,
			al.title AS album_title,
			u.username, u.avatar_url AS user_avatar_url
		FROM room_now_playing np
		JOIN tracks t ON t.id = np.track_id
		JOIN artists a ON a.id = t.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		LEFT JOIN users u ON u.id = np.suggested_by
		WHERE np.room_id = $1
	`, roomID)
	if err != nil {
		return nil, err
	}

	duration := row.DurationSeconds
	resp := &NowPlayingResponse{
		Track: TrackSummaryResponse{
			ID:              row.TrackID.String(),
			Title:           row.Title,
			ArtistName:      row.ArtistName,
			AlbumTitle:      row.AlbumTitle,
			CoverURL:        row.CoverURL,
			DurationSeconds: &duration,
		},
		StartedAt: row.StartedAt,
		Source:    row.Source,
	}

	if row.SuggestedBy != nil && row.Username != nil {
		username := *row.Username
		resp.SuggestedBy = &UserSummaryResponse{
			ID:        row.SuggestedBy.String(),
			Username:  username,
			AvatarURL: row.UserAvatarURL,
		}
	}

	return resp, nil
}

func (r *repository) GetCandidatesWithTrackAndUser(ctx context.Context, roomID, userID uuid.UUID) ([]CandidateResponse, error) {
	var rows []candidateTrackRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT
			c.id, c.room_id, c.track_id, c.suggested_by, c.vote_count, c.created_at,
			CASE WHEN v.id IS NOT NULL THEN true ELSE false END AS has_voted,
			t.title, t.duration_seconds, t.cover_url,
			a.name AS artist_name,
			al.title AS album_title,
			u.username, u.avatar_url AS user_avatar_url
		FROM room_queue_candidates c
		JOIN tracks t ON t.id = c.track_id
		JOIN artists a ON a.id = t.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		JOIN users u ON u.id = c.suggested_by
		LEFT JOIN room_queue_votes v ON v.candidate_id = c.id AND v.user_id = $2
		WHERE c.room_id = $1
		ORDER BY c.vote_count DESC, c.created_at ASC
	`, roomID, userID)
	if err != nil {
		return nil, err
	}

	items := make([]CandidateResponse, 0, len(rows))
	for _, row := range rows {
		duration := row.DurationSeconds
		items = append(items, CandidateResponse{
			ID:     row.ID,
			RoomID: row.RoomID,
			Track: TrackSummaryResponse{
				ID:              row.TrackID.String(),
				Title:           row.Title,
				ArtistName:      row.ArtistName,
				AlbumTitle:      row.AlbumTitle,
				CoverURL:        row.CoverURL,
				DurationSeconds: &duration,
			},
			SuggestedBy: UserSummaryResponse{
				ID:        row.SuggestedBy.String(),
				Username:  row.Username,
				AvatarURL: row.UserAvatarURL,
			},
			VoteCount: row.VoteCount,
			HasVoted:  row.HasVoted,
			CreatedAt: row.CreatedAt,
		})
	}

	return items, nil
}

func (r *repository) SetNowPlaying(ctx context.Context, np *RoomNowPlaying) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Upsert room_now_playing
	if _, err := tx.NamedExecContext(ctx, `
		INSERT INTO room_now_playing (room_id, track_id, started_at, suggested_by, source)
		VALUES (:room_id, :track_id, :started_at, :suggested_by, :source)
		ON CONFLICT (room_id) DO UPDATE SET
			track_id = EXCLUDED.track_id,
			started_at = EXCLUDED.started_at,
			suggested_by = EXCLUDED.suggested_by,
			source = EXCLUDED.source
	`, np); err != nil {
		return err
	}

	// If this room corresponds to a listening party, sync its current_track_id
	if _, err := tx.ExecContext(ctx, `
		UPDATE listening_parties
		SET current_track_id = $2
		WHERE id = $1 AND status != 'ended'
	`, np.RoomID, np.TrackID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) RemoveCandidate(ctx context.Context, candidateID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM room_queue_candidates WHERE id = $1
	`, candidateID)
	return err
}

func (r *repository) GetCandidatesWithVoteState(ctx context.Context, roomID, userID uuid.UUID) ([]CandidateWithVoteState, error) {
	items := make([]CandidateWithVoteState, 0)
	err := r.db.SelectContext(ctx, &items, `
		SELECT
			c.id, c.room_id, c.track_id, c.suggested_by, c.vote_count, c.created_at,
			CASE WHEN v.id IS NOT NULL THEN true ELSE false END AS has_voted
		FROM room_queue_candidates c
		LEFT JOIN room_queue_votes v ON v.candidate_id = c.id AND v.user_id = $2
		WHERE c.room_id = $1
		ORDER BY c.vote_count DESC, c.created_at ASC
	`, roomID, userID)
	return items, err
}

func (r *repository) LockRoomQueue(ctx context.Context, roomID uuid.UUID) (func(), error) {
	_, err := r.db.ExecContext(ctx, `
		SELECT pg_advisory_xact_lock(hashtext($1))
	`, "room_queue:"+roomID.String())
	if err != nil {
		return nil, err
	}
	return func() {}, nil
}

func (r *repository) PickRandomTrack(ctx context.Context) (*uuid.UUID, error) {
	var trackID uuid.UUID
	err := r.db.GetContext(ctx, &trackID, `
		SELECT id FROM tracks
		WHERE is_active = true
		ORDER BY RANDOM()
		LIMIT 1
	`)
	if err != nil {
		return nil, err
	}
	return &trackID, nil
}

func (r *repository) GetRoomMembers(ctx context.Context, roomID uuid.UUID) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `
		SELECT COUNT(*) FROM live_room_participants
		WHERE room_id = $1 AND is_active = true
	`, roomID)
	return count, err
}

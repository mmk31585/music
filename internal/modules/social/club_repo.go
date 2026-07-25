package social

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// --- Music Clubs ---

func (r *repository) CreateClub(ctx context.Context, club *MusicClub) error {
	club.ID = uuid.New()
	club.MemberCount = 1
	club.CreatedAt = time.Now().UTC()
	club.UpdatedAt = club.CreatedAt
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.NamedExecContext(ctx, `
		INSERT INTO music_clubs (id, name, slug, description, cover_url, genre, playlist_id, created_by, is_public, max_members, member_count, created_at, updated_at)
		VALUES (:id, :name, :slug, :description, :cover_url, :genre, :playlist_id, :created_by, :is_public, :max_members, :member_count, :created_at, :updated_at)
	`, club); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO music_club_members (club_id, user_id, role)
		VALUES ($1, $2, 'admin')
	`, club.ID, club.CreatedBy); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) GetClub(ctx context.Context, id uuid.UUID) (*MusicClub, error) {
	var club MusicClub
	err := r.db.GetContext(ctx, &club, `SELECT * FROM music_clubs WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &club, nil
}

func (r *repository) ListClubs(ctx context.Context, limit, offset int) ([]MusicClub, error) {
	var items []MusicClub
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM music_clubs
		WHERE is_public = true
		ORDER BY member_count DESC, created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	return items, err
}

func (r *repository) ListClubsByGenre(ctx context.Context, genre string, limit, offset int) ([]MusicClub, error) {
	var items []MusicClub
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM music_clubs
		WHERE is_public = true AND genre = $1
		ORDER BY member_count DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`, genre, limit, offset)
	return items, err
}

func (r *repository) UpdateClubPlaylistID(ctx context.Context, clubID, playlistID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE music_clubs SET playlist_id = $2, updated_at = NOW() WHERE id = $1
	`, clubID, playlistID)
	return err
}

func (r *repository) JoinClub(ctx context.Context, clubID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO music_club_members (club_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, clubID, userID)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
		UPDATE music_clubs SET member_count = (SELECT COUNT(*) FROM music_club_members WHERE club_id = $1), updated_at = NOW()
		WHERE id = $1
	`, clubID)
	return err
}

func (r *repository) LeaveClub(ctx context.Context, clubID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM music_club_members WHERE club_id = $1 AND user_id = $2
	`, clubID, userID)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
		UPDATE music_clubs SET member_count = (SELECT COUNT(*) FROM music_club_members WHERE club_id = $1), updated_at = NOW()
		WHERE id = $1
	`, clubID)
	return err
}

func (r *repository) IsClubMember(ctx context.Context, clubID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (SELECT 1 FROM music_club_members WHERE club_id = $1 AND user_id = $2)
	`, clubID, userID)
	return exists, err
}

func (r *repository) IsClubAdmin(ctx context.Context, clubID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (SELECT 1 FROM music_club_members WHERE club_id = $1 AND user_id = $2 AND role IN ('admin', 'moderator'))
	`, clubID, userID)
	return exists, err
}

func (r *repository) GetClubMembers(ctx context.Context, clubID uuid.UUID) ([]MusicClubMember, error) {
	var items []MusicClubMember
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM music_club_members WHERE club_id = $1 ORDER BY role ASC, joined_at ASC
	`, clubID)
	return items, err
}

func (r *repository) CreateClubPost(ctx context.Context, post *MusicClubPost) error {
	post.ID = uuid.New()
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO music_club_posts (id, club_id, user_id, content, created_at, updated_at)
		VALUES (:id, :club_id, :user_id, :content, :created_at, :updated_at)
	`, post)
	return err
}

func (r *repository) GetClubPosts(ctx context.Context, clubID uuid.UUID, limit, offset int) ([]MusicClubPost, error) {
	var items []MusicClubPost
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM music_club_posts WHERE club_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, clubID, limit, offset)
	return items, err
}

// --- Discussions ---

func (r *repository) CreateDiscussion(ctx context.Context, d *Discussion) error {
	d.ID = uuid.New()
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO discussions (id, user_id, target_type, target_id, content, parent_id, created_at, updated_at)
		VALUES (:id, :user_id, :target_type, :target_id, :content, :parent_id, :created_at, :updated_at)
	`, d)
	return err
}

func (r *repository) GetDiscussions(ctx context.Context, targetType string, targetID uuid.UUID, limit, offset int) ([]Discussion, error) {
	var items []Discussion
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM discussions
		WHERE target_type = $1 AND target_id = $2 AND parent_id IS NULL
		ORDER BY created_at DESC LIMIT $3 OFFSET $4
	`, targetType, targetID, limit, offset)
	return items, err
}

func (r *repository) GetDiscussionReplies(ctx context.Context, parentID uuid.UUID) ([]Discussion, error) {
	var items []Discussion
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM discussions WHERE parent_id = $1 ORDER BY created_at ASC
	`, parentID)
	return items, err
}

// --- Club Discussions (Phase 6) ---

func (r *repository) CreateClubDiscussion(ctx context.Context, d *ClubDiscussion) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO club_discussions (id, club_id, author_id, title, body, reply_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, d.ID, d.ClubID, d.AuthorID, d.Title, d.Body, d.ReplyCount, d.CreatedAt, d.UpdatedAt)
	return err
}

func (r *repository) ListClubDiscussions(ctx context.Context, clubID uuid.UUID, limit, offset int) ([]ClubDiscussion, error) {
	var items []ClubDiscussion
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, club_id, author_id, title, body, reply_count, created_at, updated_at
		FROM club_discussions
		WHERE club_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, clubID, limit, offset)
	return items, err
}

func (r *repository) GetClubDiscussion(ctx context.Context, id uuid.UUID) (*ClubDiscussion, error) {
	var d ClubDiscussion
	err := r.db.GetContext(ctx, &d, `
		SELECT id, club_id, author_id, title, body, reply_count, created_at, updated_at
		FROM club_discussions
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *repository) DeleteClubDiscussion(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM club_discussions WHERE id = $1`, id)
	return err
}

func (r *repository) IncrementClubDiscussionReplyCount(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE club_discussions SET reply_count = reply_count + 1 WHERE id = $1
	`, id)
	return err
}

func (r *repository) DecrementClubDiscussionReplyCount(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE club_discussions SET reply_count = GREATEST(reply_count - 1, 0) WHERE id = $1
	`, id)
	return err
}

func (r *repository) CreateClubDiscussionReply(ctx context.Context, reply *ClubDiscussionReply) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO club_discussion_replies (id, discussion_id, author_id, body, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, reply.ID, reply.DiscussionID, reply.AuthorID, reply.Body, reply.CreatedAt)
	return err
}

func (r *repository) GetClubDiscussionReplies(ctx context.Context, discussionID uuid.UUID) ([]ClubDiscussionReply, error) {
	var items []ClubDiscussionReply
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, discussion_id, author_id, body, created_at
		FROM club_discussion_replies
		WHERE discussion_id = $1
		ORDER BY created_at ASC
	`, discussionID)
	return items, err
}

func (r *repository) GetClubDiscussionReply(ctx context.Context, id uuid.UUID) (*ClubDiscussionReply, error) {
	var reply ClubDiscussionReply
	err := r.db.GetContext(ctx, &reply, `
		SELECT id, discussion_id, author_id, body, created_at
		FROM club_discussion_replies
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func (r *repository) DeleteClubDiscussionReply(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM club_discussion_replies WHERE id = $1`, id)
	return err
}

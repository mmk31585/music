package social

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDiscussionNotFound = fmt.Errorf("discussion not found")
	ErrReplyNotFound      = fmt.Errorf("reply not found")
	ErrNotAuthorOrOwner   = fmt.Errorf("only the author or a club admin/moderator can perform this action")
)

type DiscussionService struct {
	repo Repository
}

func NewDiscussionService(repo Repository) *DiscussionService {
	return &DiscussionService{repo: repo}
}

func (ds *DiscussionService) CreateDiscussion(ctx context.Context, clubID, authorID, title, body string) (*ClubDiscussion, error) {
	cid, err := uuid.Parse(clubID)
	if err != nil {
		return nil, fmt.Errorf("invalid club id: %w", err)
	}
	uid, err := uuid.Parse(authorID)
	if err != nil {
		return nil, fmt.Errorf("invalid author id: %w", err)
	}

	isMember, err := ds.repo.IsClubMember(ctx, cid, uid)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrNotClubMember
	}

	now := time.Now().UTC()
	d := &ClubDiscussion{
		ID:         uuid.New().String(),
		ClubID:     clubID,
		AuthorID:   authorID,
		Title:      title,
		Body:       body,
		ReplyCount: 0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := ds.repo.CreateClubDiscussion(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (ds *DiscussionService) ListDiscussions(ctx context.Context, clubID string, limit, offset int) ([]ClubDiscussion, error) {
	cid, err := uuid.Parse(clubID)
	if err != nil {
		return nil, fmt.Errorf("invalid club id: %w", err)
	}
	return ds.repo.ListClubDiscussions(ctx, cid, limit, offset)
}

func (ds *DiscussionService) GetDiscussion(ctx context.Context, discussionID string) (*ClubDiscussion, error) {
	did, err := uuid.Parse(discussionID)
	if err != nil {
		return nil, fmt.Errorf("invalid discussion id: %w", err)
	}
	d, err := ds.repo.GetClubDiscussion(ctx, did)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDiscussionNotFound
		}
		return nil, err
	}
	return d, nil
}

func (ds *DiscussionService) CreateReply(ctx context.Context, discussionID, authorID, body string) (*ClubDiscussionReply, error) {
	did, err := uuid.Parse(discussionID)
	if err != nil {
		return nil, fmt.Errorf("invalid discussion id: %w", err)
	}
	uid, err := uuid.Parse(authorID)
	if err != nil {
		return nil, fmt.Errorf("invalid author id: %w", err)
	}

	d, err := ds.repo.GetClubDiscussion(ctx, did)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDiscussionNotFound
		}
		return nil, err
	}

	cid, err := uuid.Parse(d.ClubID)
	if err != nil {
		return nil, fmt.Errorf("invalid club id: %w", err)
	}

	isMember, err := ds.repo.IsClubMember(ctx, cid, uid)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrNotClubMember
	}

	now := time.Now().UTC()
	reply := &ClubDiscussionReply{
		ID:           uuid.New().String(),
		DiscussionID: discussionID,
		AuthorID:     authorID,
		Body:         body,
		CreatedAt:    now,
	}

	if err := ds.repo.CreateClubDiscussionReply(ctx, reply); err != nil {
		return nil, err
	}

	_ = ds.repo.IncrementClubDiscussionReplyCount(ctx, did)

	return reply, nil
}

func (ds *DiscussionService) GetReplies(ctx context.Context, discussionID string) ([]ClubDiscussionReply, error) {
	did, err := uuid.Parse(discussionID)
	if err != nil {
		return nil, fmt.Errorf("invalid discussion id: %w", err)
	}
	return ds.repo.GetClubDiscussionReplies(ctx, did)
}

func (ds *DiscussionService) DeleteDiscussion(ctx context.Context, discussionID, userID string) error {
	did, err := uuid.Parse(discussionID)
	if err != nil {
		return fmt.Errorf("invalid discussion id: %w", err)
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	d, err := ds.repo.GetClubDiscussion(ctx, did)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrDiscussionNotFound
		}
		return err
	}

	cid, err := uuid.Parse(d.ClubID)
	if err != nil {
		return fmt.Errorf("invalid club id: %w", err)
	}

	isAdmin, err := ds.repo.IsClubAdmin(ctx, cid, uid)
	if err != nil {
		return err
	}

	if d.AuthorID != userID && !isAdmin {
		return ErrNotAuthorOrOwner
	}

	return ds.repo.DeleteClubDiscussion(ctx, did)
}

func (ds *DiscussionService) DeleteReply(ctx context.Context, replyID, userID string) error {
	rid, err := uuid.Parse(replyID)
	if err != nil {
		return fmt.Errorf("invalid reply id: %w", err)
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	reply, err := ds.repo.GetClubDiscussionReply(ctx, rid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrReplyNotFound
		}
		return err
	}

	did, err := uuid.Parse(reply.DiscussionID)
	if err != nil {
		return fmt.Errorf("invalid discussion id: %w", err)
	}

	d, err := ds.repo.GetClubDiscussion(ctx, did)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrDiscussionNotFound
		}
		return err
	}

	cid, err := uuid.Parse(d.ClubID)
	if err != nil {
		return fmt.Errorf("invalid club id: %w", err)
	}

	isAdmin, err := ds.repo.IsClubAdmin(ctx, cid, uid)
	if err != nil {
		return err
	}

	if reply.AuthorID != userID && !isAdmin {
		return ErrNotAuthorOrOwner
	}

	if err := ds.repo.DeleteClubDiscussionReply(ctx, rid); err != nil {
		return err
	}

	_ = ds.repo.DecrementClubDiscussionReplyCount(ctx, did)

	return nil
}

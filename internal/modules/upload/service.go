package upload

import (
	"context"
	"fmt"
	"net/http"

	apperrors "music/internal/common/errors"
	"music/internal/platform/events"

	"github.com/google/uuid"
)

type ServiceInterface interface {
	CreateDraft(ctx context.Context, userID string, req CreateUploadRequest) (Draft, error)
	GetDraft(ctx context.Context, draftID string) (Draft, error)
	ListPendingDrafts(ctx context.Context, source *string, limit, offset int) ([]Draft, int, error)
	ListUserDrafts(ctx context.Context, userID string, limit, offset int) ([]Draft, int, error)
	ReviewDraft(ctx context.Context, draftID string, reviewerID string, req ReviewDraftRequest) error
	AttachFile(ctx context.Context, draftID, userID, storageKey, fileURL string, fileSize int64) error
	AddCoUploader(ctx context.Context, trackID int64, req CoUploaderRequest) (CoUploader, error)
	RemoveCoUploader(ctx context.Context, trackID int64, userID string) error
	GetCoUploaders(ctx context.Context, trackID int64) ([]CoUploader, error)
	GetUploadSlots(ctx context.Context, userID string) (UploadSlotStatus, error)
}

type Service struct {
	repo RepositoryInterface
	bus  *events.Bus
}

func NewService(repo RepositoryInterface, bus *events.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

func (s *Service) CreateDraft(ctx context.Context, userID string, req CreateUploadRequest) (Draft, error) {
	// Check upload slots
	slots, err := s.repo.GetUploadSlots(ctx, userID)
	if err != nil {
		return Draft{}, err
	}

	if slots.UsedSlots >= slots.MaxSlots {
		return Draft{}, apperrors.New(http.StatusTooManyRequests, apperrors.CodeTooManyRequests,
			"upload slot limit reached", map[string]any{
				"used": slots.UsedSlots,
				"max":  slots.MaxSlots,
			})
	}

	// Create the draft
	draft := Draft{
		ID:                fmt.Sprintf("draft_%s_%d", userID, slots.UsedSlots+1),
		UploadedBy:        userID,
		OriginalFilename:  req.Filename,
		FilePath:          "",
		Format:            "",
		Status:            "pending",
		ExtractedMetadata: "{}",
		UploadSource:      string(SourceListener),
		ClubID:            req.ClubID,
		NeedsReview:       true,
	}

	created, err := s.repo.CreateDraft(ctx, draft)
	if err != nil {
		return Draft{}, err
	}

	// Increment slots
	if err := s.repo.IncrementSlots(ctx, userID); err != nil {
		return Draft{}, err
	}

	// Publish event
	if s.bus != nil {
		uid, _ := uuid.Parse(userID)
		s.bus.Publish(ctx, events.TrackUploadedEvent{
			BaseEvent: events.NewBaseEvent(),
			UserID:    uid,
			DraftID:   created.ID,
			Title:     req.Filename,
		})
	}

	return created, nil
}

func (s *Service) GetDraft(ctx context.Context, draftID string) (Draft, error) {
	return s.repo.GetDraft(ctx, draftID)
}

func (s *Service) ListPendingDrafts(ctx context.Context, source *string, limit, offset int) ([]Draft, int, error) {
	return s.repo.ListPendingDrafts(ctx, source, limit, offset)
}

func (s *Service) ListUserDrafts(ctx context.Context, userID string, limit, offset int) ([]Draft, int, error) {
	return s.repo.ListUserDrafts(ctx, userID, limit, offset)
}

func (s *Service) ReviewDraft(ctx context.Context, draftID string, reviewerID string, req ReviewDraftRequest) error {
	draft, err := s.repo.GetDraft(ctx, draftID)
	if err != nil {
		return err
	}

	if draft.Status != "pending" {
		return apperrors.New(http.StatusConflict, apperrors.CodeConflict, "draft is not pending review", nil)
	}

	var newStatus string
	if req.Action == "accept" {
		newStatus = "accepted"
	} else {
		newStatus = "rejected"
	}

	if err := s.repo.UpdateDraftStatus(ctx, draftID, newStatus, req.Notes, &reviewerID); err != nil {
		return err
	}

	// Publish event for accepted uploads
	if s.bus != nil && newStatus == "accepted" {
		uid, _ := uuid.Parse(draft.UploadedBy)
		s.bus.Publish(ctx, events.UploadPublishedEvent{
			BaseEvent: events.NewBaseEvent(),
			UserID:    uid,
			DraftID:   draft.ID,
			Title:     draft.OriginalFilename,
		})
	}

	// Decrement slots when draft is finalized (accepted or rejected)
	if err := s.repo.DecrementSlots(ctx, draft.UploadedBy); err != nil {
		return err
	}

	return nil
}

func (s *Service) AttachFile(ctx context.Context, draftID, userID, storageKey, fileURL string, fileSize int64) error {
	draft, err := s.repo.GetDraft(ctx, draftID)
	if err != nil {
		return err
	}
	if draft.UploadedBy != userID {
		return apperrors.New(http.StatusForbidden, apperrors.CodeForbidden, "not your draft", nil)
	}
	return s.repo.AttachFile(ctx, draftID, storageKey, fileURL, fileSize)
}

func (s *Service) AddCoUploader(ctx context.Context, trackID int64, req CoUploaderRequest) (CoUploader, error) {
	co := CoUploader{
		TrackID:           trackID,
		UserID:            req.UserID,
		Role:              req.Role,
		XpSharePercent:    req.XpSharePercent,
		ContributionType:  req.ContributionType,
	}
	return s.repo.AddCoUploader(ctx, co)
}

func (s *Service) RemoveCoUploader(ctx context.Context, trackID int64, userID string) error {
	return s.repo.RemoveCoUploader(ctx, trackID, userID)
}

func (s *Service) GetCoUploaders(ctx context.Context, trackID int64) ([]CoUploader, error) {
	return s.repo.GetCoUploaders(ctx, trackID)
}

func (s *Service) GetUploadSlots(ctx context.Context, userID string) (UploadSlotStatus, error) {
	slots, err := s.repo.GetUploadSlots(ctx, userID)
	if err != nil {
		return UploadSlotStatus{}, err
	}
	return UploadSlotStatus{
		Available: slots.UsedSlots < slots.MaxSlots,
		Used:      slots.UsedSlots,
		Max:       slots.MaxSlots,
	}, nil
}

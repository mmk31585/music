package moderation

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ReportContent(ctx context.Context, reporterID, targetID, targetType, reason string, description *string) error {
	rid, _ := uuid.Parse(reporterID)
	report := ContentReport{
		ID:          uuid.New(),
		ReporterID:  rid,
		TargetID:    targetID,
		TargetType:  targetType,
		Reason:      reason,
		Description: description,
		Status:      "pending",
		CreatedAt:   time.Now().UTC(),
	}
	return s.repo.CreateReport(ctx, report)
}

func (s *Service) GetPendingReports(ctx context.Context, limit, offset int) ([]ContentReport, error) {
	return s.repo.GetPendingReports(ctx, limit, offset)
}

func (s *Service) GetReportsByStatus(ctx context.Context, status string, limit, offset int) ([]ContentReport, error) {
	return s.repo.GetReportsByStatus(ctx, status, limit, offset)
}

func (s *Service) GetReportByID(ctx context.Context, id string) (*ContentReport, error) {
	rid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	report, err := s.repo.GetReportByID(ctx, rid)
	if err != nil {
		return nil, err
	}
	flags, _ := s.repo.GetFlags(ctx, report.TargetID, report.TargetType)
	if flags == nil {
		flags = []ContentFlag{}
	}
	return report, nil
}

func (s *Service) ResolveReport(ctx context.Context, reportID, moderatorID, status, note string) error {
	rid, _ := uuid.Parse(reportID)
	mid, _ := uuid.Parse(moderatorID)

	prev, err := s.repo.GetReportByID(ctx, rid)
	if err != nil {
		prev = nil
	}

	var prevStatus string
	if prev != nil {
		prevStatus = prev.Status
	}

	if err := s.repo.ResolveReport(ctx, rid, mid, status, note); err != nil {
		return err
	}

	action := ModerationAction{
		ID:          uuid.New(),
		ReportID:    &rid,
		ModeratorID: mid,
		Action:      status,
		TargetID:    prev.TargetID,
		TargetType:  prev.TargetType,
		NewStatus:   &status,
		Note:        &note,
		CreatedAt:   time.Now().UTC(),
	}
	if prev != nil && prevStatus != "" {
		action.PreviousStatus = &prevStatus
	}

	return s.repo.CreateAction(ctx, action)
}

func (s *Service) BulkResolve(ctx context.Context, ids []string, moderatorID, status, note string) error {
	mid, _ := uuid.Parse(moderatorID)
	uids := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		uid, err := uuid.Parse(id)
		if err == nil {
			uids = append(uids, uid)
		}
	}
	if err := s.repo.BulkResolve(ctx, uids, mid, status, note); err != nil {
		return err
	}

	for _, rid := range uids {
		report, err := s.repo.GetReportByID(ctx, rid)
		if err != nil || report == nil {
			continue
		}
		_ = s.repo.CreateAction(ctx, ModerationAction{
			ID:          uuid.New(),
			ReportID:    &rid,
			ModeratorID: mid,
			Action:      "bulk_" + status,
			TargetID:    report.TargetID,
			TargetType:  report.TargetType,
			NewStatus:   &status,
			Note:        &note,
			CreatedAt:   time.Now().UTC(),
		})
	}
	return nil
}

func (s *Service) FlagContent(ctx context.Context, targetID, targetType, flagType string, expiresInHours int) error {
	flag := ContentFlag{
		ID:         uuid.New(),
		TargetID:   targetID,
		TargetType: targetType,
		FlagType:   flagType,
		FlaggedAt:  time.Now().UTC(),
	}
	if expiresInHours > 0 {
		t := time.Now().UTC().Add(time.Duration(expiresInHours) * time.Hour)
		flag.ExpiresAt = &t
	}
	return s.repo.FlagContent(ctx, flag)
}

func (s *Service) GetAllFlags(ctx context.Context, limit, offset int, includeExpired bool) ([]ContentFlag, error) {
	return s.repo.GetAllFlags(ctx, limit, offset, includeExpired)
}

func (s *Service) GetStats(ctx context.Context) (*ModerationStats, error) {
	return s.repo.GetStats(ctx)
}

func (s *Service) GetActions(ctx context.Context, reportID string, limit, offset int) ([]ModerationAction, error) {
	var rid *uuid.UUID
	if reportID != "" {
		parsed, err := uuid.Parse(reportID)
		if err == nil {
			rid = &parsed
		}
	}
	return s.repo.GetActions(ctx, rid, limit, offset)
}

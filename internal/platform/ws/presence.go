package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type PresenceInfo struct {
	UserID    uuid.UUID `json:"user_id"`
	TrackID   *string   `json:"track_id,omitempty"`
	ArtistID  *string   `json:"artist_id,omitempty"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PresenceManager struct {
	rdb    *redis.Client
	logger *zap.Logger
	hub    *Hub
	ttl    time.Duration
}

func NewPresenceManager(rdb *redis.Client, hub *Hub, logger *zap.Logger) *PresenceManager {
	return &PresenceManager{
		rdb:    rdb,
		logger: logger,
		hub:    hub,
		ttl:    2 * time.Minute,
	}
}

func (pm *PresenceManager) UpdatePresence(ctx context.Context, userID uuid.UUID, info PresenceInfo) error {
	key := fmt.Sprintf("presence:%s", userID.String())
	info.UpdatedAt = time.Now().UTC()

	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshal presence: %w", err)
	}

	if err := pm.rdb.Set(ctx, key, string(data), pm.ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	if info.Status == "online" {
		pm.rdb.SAdd(ctx, "presence:online", userID.String())
		pm.rdb.Expire(ctx, "presence:online", pm.ttl)
	} else {
		pm.rdb.SRem(ctx, "presence:online", userID.String())
	}

	pm.hub.BroadcastToChannel("presence", Message{
		Type:    "presence.update",
		Payload: info,
	}, uuid.Nil)

	return nil
}

func (pm *PresenceManager) GetPresence(ctx context.Context, userID uuid.UUID) (*PresenceInfo, error) {
	key := fmt.Sprintf("presence:%s", userID.String())
	data, err := pm.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}

	var info PresenceInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return nil, fmt.Errorf("unmarshal presence: %w", err)
	}
	return &info, nil
}

func (pm *PresenceManager) GetOnlineFriends(ctx context.Context, followerID uuid.UUID) ([]PresenceInfo, error) {
	followedIDs, err := pm.rdb.SMembers(ctx, fmt.Sprintf("following:%s", followerID.String())).Result()
	if err != nil {
		return nil, fmt.Errorf("redis smembers: %w", err)
	}

	onlineIDs, err := pm.rdb.SMembers(ctx, "presence:online").Result()
	if err != nil {
		return nil, fmt.Errorf("redis smembers online: %w", err)
	}

	onlineSet := make(map[string]struct{}, len(onlineIDs))
	for _, id := range onlineIDs {
		onlineSet[id] = struct{}{}
	}

	var presences []PresenceInfo
	for _, followedID := range followedIDs {
		if _, online := onlineSet[followedID]; online {
			uid, err := uuid.Parse(followedID)
			if err != nil {
				continue
			}
			info, err := pm.GetPresence(ctx, uid)
			if err != nil || info == nil {
				continue
			}
			presences = append(presences, *info)
		}
	}

	return presences, nil
}

func (pm *PresenceManager) RemovePresence(ctx context.Context, userID uuid.UUID) error {
	key := fmt.Sprintf("presence:%s", userID.String())
	pm.rdb.Del(ctx, key)
	pm.rdb.SRem(ctx, "presence:online", userID.String())

	pm.hub.BroadcastToChannel("presence", Message{
		Type: "presence.offline",
		Payload: map[string]string{
			"user_id": userID.String(),
		},
	}, uuid.Nil)

	return nil
}

func (pm *PresenceManager) TrackListening(ctx context.Context, userID uuid.UUID, trackID, artistID string) error {
	info := PresenceInfo{
		UserID:  userID,
		TrackID: &trackID,
		Status:  "listening",
	}
	if artistID != "" {
		info.ArtistID = &artistID
	}
	return pm.UpdatePresence(ctx, userID, info)
}

type FriendFeedEvent struct {
	UserID    uuid.UUID       `json:"user_id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time       `json:"timestamp"`
}

func (pm *PresenceManager) PublishFeedEvent(ctx context.Context, userID uuid.UUID, eventType string, payload any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	event := FriendFeedEvent{
		UserID:    userID,
		Type:      eventType,
		Payload:   raw,
		Timestamp: time.Now().UTC(),
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if pm.rdb != nil {
		if err := pm.rdb.Publish(ctx, "feed:events", string(eventData)).Err(); err != nil {
			pm.logger.Warn("redis publish feed event failed", zap.Error(err))
		}
	}

	pm.hub.BroadcastToChannel("feed", Message{
		Type:    "feed." + eventType,
		Payload: payload,
	}, uuid.Nil)

	return nil
}

func (pm *PresenceManager) StartPresenceCleanup(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				pm.cleanupStalePresences(ctx)
			}
		}
	}()
}

func (pm *PresenceManager) cleanupStalePresences(ctx context.Context) {
	keys, err := pm.rdb.Keys(ctx, "presence:*").Result()
	if err != nil {
		pm.logger.Warn("presence cleanup keys failed", zap.Error(err))
		return
	}

	for _, key := range keys {
		ttl, err := pm.rdb.TTL(ctx, key).Result()
		if err != nil {
			continue
		}
		if ttl <= 0 {
			userID := strings.TrimPrefix(key, "presence:")
			pm.rdb.SRem(ctx, "presence:online", userID)
		}
	}
}

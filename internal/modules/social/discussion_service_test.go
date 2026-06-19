package social

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDiscussionService_CreateRequiresMembership(t *testing.T) {
	m := new(mockRepo)
	svc := NewDiscussionService(m)
	clubID := uuid.New()
	authorID := uuid.New()

	m.On("IsClubMember", mock.Anything, clubID, authorID).Return(false, nil)

	d, err := svc.CreateDiscussion(context.Background(), clubID.String(), authorID.String(), "Title", "Body")

	assert.Error(t, err)
	assert.Equal(t, ErrNotClubMember, err)
	assert.Nil(t, d)
	m.AssertExpectations(t)
}

func TestDiscussionService_CreateAsMember(t *testing.T) {
	m := new(mockRepo)
	svc := NewDiscussionService(m)
	clubID := uuid.New()
	authorID := uuid.New()

	m.On("IsClubMember", mock.Anything, clubID, authorID).Return(true, nil)
	m.On("CreateClubDiscussion", mock.Anything, mock.AnythingOfType("*social.ClubDiscussion")).Return(nil)

	d, err := svc.CreateDiscussion(context.Background(), clubID.String(), authorID.String(), "Test Title", "Test Body")

	assert.NoError(t, err)
	assert.NotNil(t, d)
	assert.Equal(t, "Test Title", d.Title)
	assert.Equal(t, "Test Body", d.Body)
	assert.Equal(t, clubID.String(), d.ClubID)
	assert.Equal(t, authorID.String(), d.AuthorID)
	assert.Equal(t, 0, d.ReplyCount)
	m.AssertExpectations(t)
}

func TestDiscussionService_ReplyIncrementsCount(t *testing.T) {
	m := new(mockRepo)
	svc := NewDiscussionService(m)
	clubID := uuid.New()
	authorID := uuid.New()
	discussionID := uuid.New()

	m.On("GetClubDiscussion", mock.Anything, discussionID).Return(&ClubDiscussion{
		ID:     discussionID.String(),
		ClubID: clubID.String(),
	}, nil)
	m.On("IsClubMember", mock.Anything, clubID, authorID).Return(true, nil)
	m.On("CreateClubDiscussionReply", mock.Anything, mock.AnythingOfType("*social.ClubDiscussionReply")).Return(nil)
	m.On("IncrementClubDiscussionReplyCount", mock.Anything, discussionID).Return(nil)

	reply, err := svc.CreateReply(context.Background(), discussionID.String(), authorID.String(), "Reply body")

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, "Reply body", reply.Body)
	assert.Equal(t, discussionID.String(), reply.DiscussionID)
	assert.Equal(t, authorID.String(), reply.AuthorID)
	m.AssertExpectations(t)
}

func TestDiscussionService_DeleteByAuthor(t *testing.T) {
	m := new(mockRepo)
	svc := NewDiscussionService(m)
	clubID := uuid.New()
	authorID := uuid.New()
	discussionID := uuid.New()

	m.On("GetClubDiscussion", mock.Anything, discussionID).Return(&ClubDiscussion{
		ID:       discussionID.String(),
		ClubID:   clubID.String(),
		AuthorID: authorID.String(),
	}, nil)
	m.On("IsClubAdmin", mock.Anything, clubID, authorID).Return(false, nil)
	m.On("DeleteClubDiscussion", mock.Anything, discussionID).Return(nil)

	err := svc.DeleteDiscussion(context.Background(), discussionID.String(), authorID.String())

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestDiscussionService_DeleteByClubOwner(t *testing.T) {
	m := new(mockRepo)
	svc := NewDiscussionService(m)
	clubID := uuid.New()
	authorID := uuid.New()
	adminID := uuid.New()
	discussionID := uuid.New()

	m.On("GetClubDiscussion", mock.Anything, discussionID).Return(&ClubDiscussion{
		ID:       discussionID.String(),
		ClubID:   clubID.String(),
		AuthorID: authorID.String(),
	}, nil)
	m.On("IsClubAdmin", mock.Anything, clubID, adminID).Return(true, nil)
	m.On("DeleteClubDiscussion", mock.Anything, discussionID).Return(nil)

	err := svc.DeleteDiscussion(context.Background(), discussionID.String(), adminID.String())

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestDiscussionService_DeleteByNonAuthorNonOwnerFails(t *testing.T) {
	m := new(mockRepo)
	svc := NewDiscussionService(m)
	clubID := uuid.New()
	authorID := uuid.New()
	strangerID := uuid.New()
	discussionID := uuid.New()

	m.On("GetClubDiscussion", mock.Anything, discussionID).Return(&ClubDiscussion{
		ID:       discussionID.String(),
		ClubID:   clubID.String(),
		AuthorID: authorID.String(),
	}, nil)
	m.On("IsClubAdmin", mock.Anything, clubID, strangerID).Return(false, nil)

	err := svc.DeleteDiscussion(context.Background(), discussionID.String(), strangerID.String())

	assert.Error(t, err)
	assert.Equal(t, ErrNotAuthorOrOwner, err)
	m.AssertExpectations(t)
}

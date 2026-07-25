package ws

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// --- Test Helpers ---

func newTestHub(t *testing.T) *Hub {
	t.Helper()
	return NewHub(zap.NewNop(), nil)
}

func newTestClient(t *testing.T) *Client {
	t.Helper()
	return &Client{
		UserID: uuid.New(),
		Send:   make(chan []byte, 256),
	}
}

func runHub(t *testing.T, hub *Hub) context.CancelFunc {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	go hub.Run(ctx)
	t.Cleanup(cancel)
	return cancel
}

func registerClient(t *testing.T, hub *Hub, client *Client) {
	t.Helper()
	client.Hub = hub
	hub.register <- client
}

func unregisterClient(t *testing.T, hub *Hub, client *Client) {
	t.Helper()
	hub.unregister <- client
}

func assertMsgReceived(t *testing.T, client *Client, timeout time.Duration) []byte {
	t.Helper()
	select {
	case msg := <-client.Send:
		return msg
	case <-time.After(timeout):
		t.Fatalf("client %v did not receive message within %v", client.UserID, timeout)
		return nil
	}
}

func assertMsgNotReceived(t *testing.T, client *Client, timeout time.Duration) {
	t.Helper()
	select {
	case msg := <-client.Send:
		t.Fatalf("client %v unexpectedly received message: %s", client.UserID, string(msg))
	case <-time.After(timeout):
	}
}

// --- Client Registration & Unregistration ---

func TestHub_ClientRegistration(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	client := newTestClient(t)
	registerClient(t, hub, client)

	time.Sleep(10 * time.Millisecond)

	assert.True(t, hub.IsOnline(client.UserID))
	assert.Contains(t, hub.OnlineUsers(), client.UserID)
}

func TestHub_ClientUnregistration(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	client := newTestClient(t)
	registerClient(t, hub, client)
	time.Sleep(10 * time.Millisecond)

	assert.True(t, hub.IsOnline(client.UserID))

	unregisterClient(t, hub, client)
	time.Sleep(10 * time.Millisecond)

	assert.False(t, hub.IsOnline(client.UserID))
	assert.Empty(t, hub.OnlineUsers())

	_, ok := <-client.Send
	assert.False(t, ok, "client Send channel should be closed after unregistration")
}

func TestHub_UnregisterCleansSubscriptions(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	client := newTestClient(t)
	registerClient(t, hub, client)
	time.Sleep(10 * time.Millisecond)

	hub.Subscribe(client, "room:test")
	hub.Subscribe(client, "room:another")

	hub.mu.RLock()
	assert.Contains(t, hub.subscriptions["room:test"], client)
	assert.Contains(t, hub.subscriptions["room:another"], client)
	hub.mu.RUnlock()

	unregisterClient(t, hub, client)
	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	_, exists := hub.subscriptions["room:test"]
	assert.False(t, exists, "subscription should be removed after client unregisters")
	_, exists = hub.subscriptions["room:another"]
	assert.False(t, exists, "subscription should be removed after client unregisters")
	hub.mu.RUnlock()
}

func TestHub_RegisterMultipleClients(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	client1 := newTestClient(t)
	client2 := newTestClient(t)
	client3 := newTestClient(t)

	registerClient(t, hub, client1)
	registerClient(t, hub, client2)
	registerClient(t, hub, client3)
	time.Sleep(10 * time.Millisecond)

	users := hub.OnlineUsers()
	assert.Len(t, users, 3)
	assert.True(t, hub.IsOnline(client1.UserID))
	assert.True(t, hub.IsOnline(client2.UserID))
	assert.True(t, hub.IsOnline(client3.UserID))
}

// --- Channel Subscription & Unsubscription ---

func TestHub_SubscribeChannel(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	client := newTestClient(t)
	registerClient(t, hub, client)
	time.Sleep(10 * time.Millisecond)

	ok := hub.Subscribe(client, "room:abc")
	assert.True(t, ok)

	hub.mu.RLock()
	subs, exists := hub.subscriptions["room:abc"]
	hub.mu.RUnlock()
	assert.True(t, exists)
	assert.Contains(t, subs, client)
}

func TestHub_UnsubscribeChannel(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	client := newTestClient(t)
	registerClient(t, hub, client)
	time.Sleep(10 * time.Millisecond)

	hub.Subscribe(client, "room:abc")
	hub.Unsubscribe(client, "room:abc")

	hub.mu.RLock()
	_, exists := hub.subscriptions["room:abc"]
	hub.mu.RUnlock()
	assert.False(t, exists, "channel should be deleted when no subscribers remain")
}

func TestHub_SubscribeWithAuthorization(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	hub.Authorize = func(userID uuid.UUID, channel string) bool {
		return channel == "allowed-channel"
	}

	client := newTestClient(t)
	registerClient(t, hub, client)
	time.Sleep(10 * time.Millisecond)

	assert.True(t, hub.Subscribe(client, "allowed-channel"))
	assert.False(t, hub.Subscribe(client, "denied-channel"))

	hub.mu.RLock()
	_, allowedExists := hub.subscriptions["allowed-channel"]
	_, deniedExists := hub.subscriptions["denied-channel"]
	hub.mu.RUnlock()
	assert.True(t, allowedExists)
	assert.False(t, deniedExists)
}

// --- Local Broadcast ---

func TestHub_BroadcastToChannel(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	client1 := newTestClient(t)
	client2 := newTestClient(t)
	registerClient(t, hub, client1)
	registerClient(t, hub, client2)
	time.Sleep(10 * time.Millisecond)

	hub.Subscribe(client1, "room:chat")
	hub.Subscribe(client2, "room:chat")

	msg := Message{Type: "chat.message", Payload: "hello"}
	hub.BroadcastToChannel("room:chat", msg, uuid.Nil)

	got1 := assertMsgReceived(t, client1, 100*time.Millisecond)
	got2 := assertMsgReceived(t, client2, 100*time.Millisecond)

	var decoded1, decoded2 Message
	json.Unmarshal(got1, &decoded1)
	json.Unmarshal(got2, &decoded2)
	assert.Equal(t, "chat.message", decoded1.Type)
	assert.Equal(t, "chat.message", decoded2.Type)
}

func TestHub_BroadcastExcludesSender(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	client1 := newTestClient(t)
	client2 := newTestClient(t)
	registerClient(t, hub, client1)
	registerClient(t, hub, client2)
	time.Sleep(10 * time.Millisecond)

	hub.Subscribe(client1, "room:chat")
	hub.Subscribe(client2, "room:chat")

	msg := Message{Type: "chat.message", Payload: "exclude me"}
	hub.BroadcastToChannel("room:chat", msg, client1.UserID)

	assertMsgNotReceived(t, client1, 50*time.Millisecond)
	assertMsgReceived(t, client2, 100*time.Millisecond)
}

func TestHub_BroadcastToUnsubscribedChannel(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	client := newTestClient(t)
	registerClient(t, hub, client)
	time.Sleep(10 * time.Millisecond)

	msg := Message{Type: "test", Payload: "no one"}
	hub.BroadcastToChannel("empty-channel", msg, uuid.Nil)

	assertMsgNotReceived(t, client, 50*time.Millisecond)
}

func TestHub_BroadcastMultipleChannels(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	client := newTestClient(t)
	registerClient(t, hub, client)
	time.Sleep(10 * time.Millisecond)

	hub.Subscribe(client, "channel:a")
	hub.Subscribe(client, "channel:b")

	msgA := Message{Type: "a", Payload: "only a"}
	msgB := Message{Type: "b", Payload: "only b"}

	hub.BroadcastToChannel("channel:a", msgA, uuid.Nil)
	hub.BroadcastToChannel("channel:b", msgB, uuid.Nil)

	gotA := assertMsgReceived(t, client, 100*time.Millisecond)
	gotB := assertMsgReceived(t, client, 100*time.Millisecond)

	// Messages may arrive in any order; just verify both types arrive
	var decoded Message
	var types [2]string
	json.Unmarshal(gotA, &decoded)
	types[0] = decoded.Type
	json.Unmarshal(gotB, &decoded)
	types[1] = decoded.Type

	assert.Contains(t, types[:], "a")
	assert.Contains(t, types[:], "b")
}

// --- SendToUser ---

func TestHub_SendToUser(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	client := newTestClient(t)
	registerClient(t, hub, client)
	time.Sleep(10 * time.Millisecond)

	msg := Message{Type: "direct", Payload: "private message"}
	hub.SendToUser(client.UserID, msg)

	got := assertMsgReceived(t, client, 100*time.Millisecond)
	var decoded Message
	json.Unmarshal(got, &decoded)
	assert.Equal(t, "direct", decoded.Type)
}

func TestHub_SendToNonExistentUser(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	msg := Message{Type: "direct", Payload: "nobody"}
	hub.SendToUser(uuid.New(), msg)
}

func TestHub_SendToUserUnregistered(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	msg := Message{Type: "direct", Payload: "ghost"}
	hub.SendToUser(uuid.New(), msg)
}

// --- OnlineUsers / IsOnline ---

func TestHub_OnlineUsers(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	assert.Empty(t, hub.OnlineUsers())

	client1 := newTestClient(t)
	client2 := newTestClient(t)

	registerClient(t, hub, client1)
	registerClient(t, hub, client2)
	time.Sleep(10 * time.Millisecond)

	users := hub.OnlineUsers()
	assert.Len(t, users, 2)
	assert.Contains(t, users, client1.UserID)
	assert.Contains(t, users, client2.UserID)

	unregisterClient(t, hub, client1)
	time.Sleep(10 * time.Millisecond)

	users = hub.OnlineUsers()
	assert.Len(t, users, 1)
	assert.Equal(t, client2.UserID, users[0])
}

func TestHub_IsOnline(t *testing.T) {
	hub := newTestHub(t)
	runHub(t, hub)

	client := newTestClient(t)
	registerClient(t, hub, client)
	time.Sleep(10 * time.Millisecond)

	assert.True(t, hub.IsOnline(client.UserID))
	assert.False(t, hub.IsOnline(uuid.New()))

	unregisterClient(t, hub, client)
	time.Sleep(10 * time.Millisecond)

	assert.False(t, hub.IsOnline(client.UserID))
}

// --- Shutdown ---

func TestHub_Shutdown(t *testing.T) {
	hub := newTestHub(t)
	ctx, cancel := context.WithCancel(context.Background())
	go hub.Run(ctx)

	client := newTestClient(t)
	registerClient(t, hub, client)
	time.Sleep(10 * time.Millisecond)

	assert.True(t, hub.IsOnline(client.UserID))

	cancel()

	time.Sleep(50 * time.Millisecond)

	_, ok := <-client.Send
	assert.False(t, ok, "client Send channel should be closed after shutdown")
}

func TestHub_ShutdownWithMultipleClients(t *testing.T) {
	hub := newTestHub(t)
	ctx, cancel := context.WithCancel(context.Background())
	go hub.Run(ctx)

	clients := make([]*Client, 5)
	for i := 0; i < 5; i++ {
		client := newTestClient(t)
		registerClient(t, hub, client)
		clients[i] = client
	}
	time.Sleep(10 * time.Millisecond)

	cancel()
	time.Sleep(50 * time.Millisecond)

	for _, client := range clients {
		_, ok := <-client.Send
		assert.False(t, ok, "all client Send channels should be closed after shutdown")
	}
}

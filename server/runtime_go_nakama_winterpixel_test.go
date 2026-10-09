package server

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestOnlineUserCountCountsUsersNotSockets(t *testing.T) {
	logger := loggerForTest(t)
	tracker := StartLocalTracker(logger, NewConfig(logger), NewLocalSessionRegistry(&testMetrics{}), nil, &testMetrics{}, &protojson.MarshalOptions{})
	defer tracker.Stop()
	nk := &RuntimeGoNakamaModule{tracker: tracker}

	// What the socket acceptor tracks for every session it accepts.
	connect := func(userID uuid.UUID) uuid.UUID {
		sessionID := uuid.Must(uuid.NewV4())
		tracker.Track(context.Background(), sessionID, PresenceStream{Mode: StreamModeNotifications, Subject: userID}, userID, PresenceMeta{Hidden: true}, true)
		return sessionID
	}

	if got := nk.OnlineUserCount(); got != 0 {
		t.Fatalf("OnlineUserCount() = %d with nobody connected, want 0", got)
	}

	userA, userB := uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4())
	firstTab := connect(userA)
	connect(userA)
	sessionB := connect(userB)

	// A presence on any other stream is not another connection.
	tracker.Track(context.Background(), sessionB, PresenceStream{Mode: StreamModeChannel, Subject: uuid.Must(uuid.NewV4())}, userB, PresenceMeta{Hidden: true}, false)

	if got := nk.OnlineUserCount(); got != 2 {
		t.Fatalf("OnlineUserCount() = %d, want 2: two users on three sockets", got)
	}

	tracker.UntrackAll(firstTab, runtime.PresenceReasonDisconnect)
	if got := nk.OnlineUserCount(); got != 2 {
		t.Errorf("OnlineUserCount() = %d after one of two tabs closed, want 2", got)
	}

	tracker.UntrackAll(sessionB, runtime.PresenceReasonDisconnect)
	if got := nk.OnlineUserCount(); got != 1 {
		t.Errorf("OnlineUserCount() = %d after a user disconnected, want 1", got)
	}
}

package server

// Winterpixel additions to the Go runtime module.
//
// These are methods on *RuntimeGoNakamaModule that are deliberately left off
// nakama-common's runtime.NakamaModule interface. Adding to that interface would
// mean forking nakama-common too, and a Go plugin only loads against the exact
// nakama-common the server was built with, so the interface stays upstream's
// and a module reaches these with a type assertion instead:
//
//	counter, ok := nk.(interface{ OnlineUserCount() int })
//
// The name and signature are the whole contract. Changing one here means
// changing the module's interface to match, or the assertion quietly fails.
// They also live in their own file so an upstream merge never touches them.

// Pinned so that a change to the signature breaks this build, rather than the
// module's type assertion at runtime.
var _ interface{ OnlineUserCount() int } = (*RuntimeGoNakamaModule)(nil)

// OnlineUserCount is how many distinct users have a live socket on this node.
//
// Every socket session tracks a hidden presence on its user's notification
// stream, and that stream is keyed by user, so there is one stream per user
// online however many sockets each of them has open. Counting the streams is
// the same walk NotificationSendAll uses to reach everyone connected. The
// tracker drops a stream when its last presence leaves, so none are empty.
//
// Node-local, like the tracker it reads.
func (n *RuntimeGoNakamaModule) OnlineUserCount() int {
	notificationStreamMode := StreamModeNotifications
	return len(n.tracker.CountByStreamModeFilter(map[uint8]*uint8{StreamModeNotifications: &notificationStreamMode}))
}

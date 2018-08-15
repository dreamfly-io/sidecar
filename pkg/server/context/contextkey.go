package context

// ContextKey type
type Key string

// Context key types
const (
	ContextKeyStreamID               Key = "StreamId"
	ContextKeyConnectionID           Key = "ConnectionId"
	ContextKeyListenerIp             Key = "ListenerIp"
	ContextKeyListenerPort           Key = "ListenerPort"
	ContextKeyListenerName           Key = "ListenerName"
	ContextKeyLogger                 Key = "Logger"
	ContextKeyAccessLogs             Key = "AccessLogs"
)

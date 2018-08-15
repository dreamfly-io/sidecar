package connection

// ConnectionEvent
type Event string

// ConnectionEvents
const (
	RemoteClose     Event = "RemoteClose"
	LocalClose      Event = "LocalClose"
	OnReadErrClose  Event = "OnReadErrClose"
	OnWriteErrClose Event = "OnWriteErrClose"
	OnConnect       Event = "OnConnect"
	Connected       Event = "ConnectedFlag"
	ConnectTimeout  Event = "ConnectTimeout"
	ConnectFailed   Event = "ConnectFailed"
)

// IsClose represents whether the event is triggered by connection close
func (e Event) IsClose() bool {
	return e == LocalClose || e == RemoteClose ||
		e == OnReadErrClose || e == OnWriteErrClose
}

// IsConnectFailure represents whether the event is triggered by connection failure
func (e Event) IsFailure() bool {
	return e == ConnectFailed || e == ConnectTimeout
}

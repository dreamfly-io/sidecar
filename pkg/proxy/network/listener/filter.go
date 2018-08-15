package listener

// FilterStatus type
type Status string

// FilterStatus types
const (
	Continue      Status = "Continue"
	StopIteration Status = "StopIteration"
)

type Filter interface {
	// OnAccept is called when a raw connection is accepted, but before a Connection is created.
	OnAccept() Status
}
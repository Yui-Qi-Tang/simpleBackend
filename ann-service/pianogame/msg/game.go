package msg

// Let front-end follow?
const (
	SendPianoKey = iota
	ExitConn     // necessary?
	ResetConn
	WelcomeConn // necessary?
)

// Base message base
type Base struct {
	To     interface{}
	From   interface{}
	Action int
}

// PianoKey for client exchange the piano key who  is pressed
type PianoKey struct {
	Base
	Key interface{}
}

// Text for client exchange the text msg
type Text struct {
	Base
	Text string
}

// Welcome for client receivce welcome msg -> a init phase finish in client side
type Welcome struct {
	ID   interface{}
	Text string
}

// Exit for client exits the game
type Exit struct {
	Base
	Text interface{}
}

// Error error message for client
type Error struct {
	Text string
}

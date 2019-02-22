package msg

// Let front-end follow?
const (
	SendPianoKey = iota
	ExitConn
	ResetConn
)

type MsgBase struct {
	To     interface{}
	From   interface{}
	Action int
}

// PianoKey for client exchange the piano key who  is pressed
type PianoKey struct {
	MsgBase
	Key interface{}
}

// Text for client exchange the text msg
type Text struct {
	MsgBase
	Text string
}

// Welcome for client receivce welcome msg -> a init phase finish in client side
type Welcome struct {
	ID   interface{}
	Text string
}

// Exit for client exits the game
type Exit struct {
	MsgBase
	Text interface{}
}

// Error error message for client
type Error struct {
	Text string
}

package msg

type msgSrcAndDest struct {
	To   interface{}
	From interface{}
}

// PianoKey for client exchange the piano key who  is pressed
type PianoKey struct {
	msgSrcAndDest
	Key interface{}
}

// Text for client exchange the text msg
type Text struct {
	msgSrcAndDest
	Text string
}

// Welcome for client receivce welcome msg -> a init phase finish in client side
type Welcome struct {
	ID   interface{}
	Text string
}

// Error error message for client
type Error struct {
	Text string
}

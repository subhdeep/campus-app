package models

// WebRTCInitMessage is sent by a client wanting to start a connection
type WebRTCInitMessage struct {
	To     Username `json:"to"`
	From   Username `json:"from"`
	FromID ConnID   `json:"fromID"`
}

// WebRTCCancelMessage is sent by a client to cancel a call
type WebRTCCancelMessage struct {
	To Username `json:"to"`
}

// WebRTCAckMessage is a response from an interested client
type WebRTCAckMessage struct {
	From   Username `json:"from"`
	FromID ConnID   `json:"fromID"`
	ToID   ConnID   `json:"toID"`
}

// WebRTCMessage is a message between two clients who have decided to
// link up via WebRTC.
type WebRTCMessage struct {
	ToID   ConnID      `json:"toID"`
	FromID ConnID      `json:"fromID"`
	Body   interface{} `json:"body"`
}

// WebRTCRejectMessage is a message sent to a caller to inform them that the
// call has been rejected.
type WebRTCRejectMessage struct {
	ToID   ConnID   `json:"toID"`
	From   Username `json:"from"`
	FromID ConnID   `json:"fromID"`
}

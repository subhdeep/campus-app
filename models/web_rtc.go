package models

// WebRTCInit is sent by a client wanting to start a connection
type WebRTCInitMessage struct {
	To     Username `json:"to"`
	From   Username `json:"from"`
	FromID ConnID   `json:"fromID"`
}

// WebRTCAck is a response from an interested client
type WebRTCAckMessage struct {
	From   Username `json:"from"`
	FromID ConnID   `json:"fromID"`
	ToID   ConnID   `json:"toID"`
}

// WebRTC is a message between two clients who have decided to
// link up via WebRTC.
type WebRTCMessage struct {
	ToID   ConnID `json:"toID"`
	FromID ConnID `json:"fromID"`
}

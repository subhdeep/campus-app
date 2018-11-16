package models

// WebRTCInit is sent by a client wanting to start a connection
type WebRTCInit struct {
	To     Username `json:"to"`
	From   Username `json:"from"`
	FromID ConnID   `json:"fromID"`
}

// WebRTCAck is a response from an interested client
type WebRTCAck struct {
	From   Username `json:"from"`
	FromID ConnID   `json:"fromID"`
	ToID   ConnID   `json:"toID"`
}

// WebRTC is a message between two clients who have decided to
// link up via WebRTC.
type WebRTC struct {
	ToID   ConnID `json:"toID"`
	FromID ConnID `json:"fromID"`
}

package types

// Device represents device
type Device struct {
	// ID is in UUID format. ref: https://datatracker.ietf.org/doc/html/rfc4122
	ID    string `json:"id"`
	Name  string `json:"name"`
	Brand string `json:"brand"`
	Model string `json:"model"`
}

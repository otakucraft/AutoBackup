package utils

// Result generic structure for API json responses.
type Result struct {
	Code uint `json:"code"`
	Data string `json:"data"`
}

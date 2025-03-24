package dto

type Event struct {
	Type string `json:"type"`
	X    int    `json:"x,omitempty"`
	Y    int    `json:"y,omitempty"`
	Key  int    `json:"key,omitempty"`
}

type ErrorResponse struct {
	ErrorMsg string `json:"error_msg"`
}

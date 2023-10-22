package common

// Health is a JSON response for healthcheck request
type Health struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// StatusMessage is a generic JSON response containing the status code and a status message.
type StatusMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

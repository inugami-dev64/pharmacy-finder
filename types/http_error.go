package types

type HttpError struct {
	StatusCode int    `json:"code"`
	Timestamp  Time   `json:"ts"`
	Message    string `json:"msg"`
}

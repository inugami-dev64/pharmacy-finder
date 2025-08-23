package types

import "time"

type HttpError struct {
	StatusCode int    `json:"code"`
	Timestamp  Time   `json:"ts"`
	Message    string `json:"msg"`
}

func NewHttpError(code int, message string) HttpError {
	return HttpError{
		StatusCode: code,
		Timestamp:  Time(time.Now().UTC()),
		Message:    message,
	}
}

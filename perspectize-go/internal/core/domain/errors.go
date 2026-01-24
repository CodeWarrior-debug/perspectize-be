package domain

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrInvalidURL    = errors.New("invalid URL")
	ErrYouTubeAPI    = errors.New("youtube API error")
)

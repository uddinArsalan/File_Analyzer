package services

import "errors"

var (
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrUserAlreadyExists  = errors.New("User already exists")
	ErrSessionExpired     = errors.New("Session expired")
	ErrDocumentNotFound   = errors.New("Document not found")
	ErrDocumentForbidden  = errors.New("Document forbidden")
	ErrAlreadyProcessing  = errors.New("Document already processing")
)

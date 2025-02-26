package model

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrInsufficientQuota = errors.New("insufficient quota")
)

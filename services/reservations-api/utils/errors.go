package utils

import (
	"errors"
)

var (
	ErrReservationNotFound    = errors.New("reservation not found")
	ErrInvalidReservationData = errors.New("invalid reservation data")
	ErrInternalServer         = errors.New("internal server error")
)

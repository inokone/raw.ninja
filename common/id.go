package common

import (
	"github.com/google/uuid"
)

func NewId() string {
	return uuid.NewString()
}

package definitions

import (
	"time"

	"github.com/google/uuid"
)

type Screening struct {
	ID        *uuid.UUID    `json:"id"`
	StartTime time.Time     `json:"startTime"`
	Duration  time.Duration `json:"duration"`
	Title     string        `json:"title"`
}

package arpio

import (
	"time"
)

type RecoveryPoint struct {
	AvailableAt     time.Time `json:"availableAt"`
	Protected       bool      `json:"protected"`
	RecoveryPointID string    `json:"recoveryPointId"`
	Timestamp       time.Time `json:"timestamp"`
}

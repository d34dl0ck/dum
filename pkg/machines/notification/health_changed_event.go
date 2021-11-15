package notification

import "github.com/google/uuid"

type HealthChangedEvent struct {
	Level     int
	MachineId uuid.UUID
}

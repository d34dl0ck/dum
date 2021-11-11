package adapters

import (
	"dum/internal/machines/entities"
	"log"
)

// Simple notification strategy writes to log
type LogNotificationStrategy struct {
	logger *log.Logger
}

func (s LogNotificationStrategy) Notify(id entities.MachineId, level entities.HealthLevel) error {
	s.logger.Printf(template, id, level)
	return nil
}

func NewLogNotificationStrategy(logger *log.Logger) entities.HealthNotificationStrategy {
	return &LogNotificationStrategy{
		logger: logger,
	}
}

const template string = "Machine %s reporting about it's health change state %d"

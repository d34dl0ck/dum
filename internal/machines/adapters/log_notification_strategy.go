package adapters

import (
	"dum/internal/machines/entities"
	"log"
)

// Simple notification strategy writes to log
type LogNotificationStrategy struct {
	logger *log.Logger
}

func (s LogNotificationStrategy) Notify(machineName string, level entities.HealthLevel) error {
	s.logger.Printf(template, machineName, level)
	return nil
}

func NewLogNotificationStrategy(logger *log.Logger) entities.HealthNotificationStrategy {
	return &LogNotificationStrategy{
		logger: logger,
	}
}

const template string = "Machine %s reporting about it's health change state %d"

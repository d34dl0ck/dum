package entities

// Interface for notifiying about machine health level changes.
type HealthNotificationStrategy interface {
	Notify(machineName string, level HealthLevel) error
}

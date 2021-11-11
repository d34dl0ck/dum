package entities

// Interface for notifiying about machine health level changes.
type HealthNotificationStrategy interface {
	Notify(machineId MachineId, level HealthLevel) error
}

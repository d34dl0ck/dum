package entities

import (
	"time"

	"github.com/google/uuid"
)

// MissingUpdate is a value type, representing the information about uninstalled update on
// some machine - update severity and duration between time update became available and
// report time.
type MissingUpdate struct {
	UpdateId uuid.UUID
	Severity Severity
	Duration time.Duration
}

// Severity represents a level of update importance.
type Severity int

const (
	Unspecified Severity = iota
	Low
	Important
	Critical
)

// Machine is an entitiy type which represents a machine
// with some health level. It can process reports about
// updates, missing on this machine.
type Machine struct {
	h       health
	Name    string
	missing []MissingUpdate
}

// Creates a machine with specific missing updates and health level.
func CreateMachine(name string, mu []MissingUpdate) *Machine {
	health := &health{}

	return &Machine{
		Name:    name,
		h:       health.Recalculate(mu),
		missing: mu,
	}
}

// Returns health level of machine.
func (m *Machine) GetHealthLevel() HealthLevel {
	return m.h.level
}

// Returns current set of missing updates
func (m *Machine) GetMissingUpdates() []MissingUpdate {
	return m.missing
}

// Processes message about missing updates appearing for this machine.
func (m *Machine) Report(mu []MissingUpdate, s HealthNotificationStrategy) error {
	m.h = m.h.Recalculate(mu)
	m.missing = mu
	err := s.Notify(m.Name, m.h.level)

	if err != nil {
		return err
	}

	return nil
}

// HealthLevel is an indicator for machine health.
type HealthLevel int

const (
	Healthy HealthLevel = iota
	Warning
	Danger
)

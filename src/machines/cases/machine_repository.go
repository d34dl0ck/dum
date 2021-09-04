package cases

import "dum/machines/entities"

// Interface for accessing and persisting machine entities. Should implement optimistic locking strategy
// to support horizontal scaling of service.
type MachineRepository interface {
	// Loading machine entity from some storage.
	Load(name string) (*entities.Machine, error)

	// Saving machine enitity to some storage.
	Save(machine *entities.Machine) error
}

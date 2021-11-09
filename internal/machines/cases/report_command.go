package cases

import "dum/internal/machines/entities"

// Command for reporting about some missing updates of some machine.
type ReportCommand struct {
	MachineName          string
	MissingUpdates       []entities.MissingUpdate
	NotificationStrategy entities.HealthNotificationStrategy
	Repository           MachineRepository
}

func (c *ReportCommand) Execute() error {
	machine, err := c.Repository.Load(c.MachineName)
	if err != nil {
		return err
	}

	if machine == nil {
		machine = entities.CreateMachine(c.MachineName, []entities.MissingUpdate{})
	}

	err = machine.Report(c.MissingUpdates, c.NotificationStrategy)
	if err != nil {
		return err
	}

	err = c.Repository.Save(machine)
	if err != nil {
		return err
	}

	return nil
}

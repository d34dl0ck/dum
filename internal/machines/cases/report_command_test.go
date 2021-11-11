package cases

import (
	"dum/internal/machines/entities"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestExecuteForNewMachine(t *testing.T) {
	strategyMock := notificationStrategyMock{}
	repositoryMock := repositoryMock{
		loadedMachine:         nil,
		shouldReturnLoadError: false,
	}

	command := ReportCommand{
		MachineName:          newMachineName,
		MissingUpdates:       expectedMissingUpdates,
		Repository:           &repositoryMock,
		NotificationStrategy: &strategyMock,
	}

	err := command.Execute()

	if err != nil {
		t.Errorf("Execute should not return error %s!", err)
	}

	if strategyMock.wasCalled == false {
		t.Errorf("Notification strategy should be called!")
	}

	if repositoryMock.savedMachine == nil {
		t.Errorf("Machine should be saved!")
	}

	if repositoryMock.savedMachine.Name != newMachineName {
		t.Errorf("Machine should be saved with expected name %s, but was %s!", newMachineName, repositoryMock.savedMachine.Name)
	}

	if len(repositoryMock.savedMachine.GetMissingUpdates()) != 0 {
		t.Errorf("Missing updated of saved machine should be empty, but has length %d", len(repositoryMock.savedMachine.GetMissingUpdates()))
	}
}

func TestExecuteForExistingMachine(t *testing.T) {
	strategyMock := notificationStrategyMock{}
	repositoryMock := repositoryMock{
		loadedMachine:         entities.CreateMachine(entities.MachineId(uuid.New()), existingMachineName, []entities.MissingUpdate{}),
		shouldReturnLoadError: false,
	}

	command := ReportCommand{
		MachineName:          existingMachineName,
		MissingUpdates:       expectedMissingUpdates,
		Repository:           &repositoryMock,
		NotificationStrategy: &strategyMock,
	}

	err := command.Execute()

	if err != nil {
		t.Errorf("Execute should not return error %s!", err)
	}

	if strategyMock.wasCalled == false {
		t.Errorf("Notification strategy should be called!")
	}

	if repositoryMock.savedMachine == nil {
		t.Errorf("Machine should be saved!")
	}

	if repositoryMock.savedMachine.Name != existingMachineName {
		t.Errorf("Machine should be saved with expected name %s, but was %s!", existingMachineName, repositoryMock.savedMachine.Name)
	}

	if len(repositoryMock.savedMachine.GetMissingUpdates()) != 0 {
		t.Errorf("Missing updated of saved machine should be empty, but has length %d", len(repositoryMock.savedMachine.GetMissingUpdates()))
	}
}

func TestExecuteReturnsLoadError(t *testing.T) {
	strategyMock := notificationStrategyMock{}
	repositoryMock := repositoryMock{
		loadedMachine:         nil,
		shouldReturnLoadError: true,
	}

	command := ReportCommand{
		MachineName:          newMachineName,
		MissingUpdates:       expectedMissingUpdates,
		Repository:           &repositoryMock,
		NotificationStrategy: &strategyMock,
	}

	err := command.Execute()

	if err == nil {
		t.Errorf("Execute should return error!")
	}

	if err != errLoad {
		t.Errorf("Error mismatch! Expected %s, but was %s!", errLoad, err)
	}

	if repositoryMock.savedMachine != nil {
		t.Errorf("Should not try to save machine if load failed!")
	}
}

func TestExecuteReturnsReportError(t *testing.T) {
	strategyMock := notificationStrategyMock{
		shouldReturnError: true,
	}
	repositoryMock := repositoryMock{
		loadedMachine:         nil,
		shouldReturnLoadError: false,
	}

	command := ReportCommand{
		MachineName:          newMachineName,
		MissingUpdates:       expectedMissingUpdates,
		Repository:           &repositoryMock,
		NotificationStrategy: &strategyMock,
	}

	err := command.Execute()

	if err == nil {
		t.Errorf("Execute should return error!")
	}

	if err != errReport {
		t.Errorf("Error mismatch! Expected %s, but was %s!", errReport, err)
	}

	if repositoryMock.savedMachine != nil {
		t.Errorf("Should not try to save machine if report failed!")
	}
}

func TestExecuteReturnsSaveError(t *testing.T) {
	strategyMock := notificationStrategyMock{
		shouldReturnError: false,
	}
	repositoryMock := repositoryMock{
		loadedMachine:         nil,
		shouldReturnLoadError: false,
		shouldReturnSaveError: true,
	}

	command := ReportCommand{
		MachineName:          newMachineName,
		MissingUpdates:       expectedMissingUpdates,
		Repository:           &repositoryMock,
		NotificationStrategy: &strategyMock,
	}

	err := command.Execute()

	if err == nil {
		t.Errorf("Execute should return error!")
	}

	if err != errSave {
		t.Errorf("Error mismatch! Expected %s, but was %s!", errReport, err)
	}

	if strategyMock.wasCalled == false {
		t.Errorf("Notification strategy should be called!")
	}
}

type notificationStrategyMock struct {
	shouldReturnError bool
	wasCalled         bool
}

func (m *notificationStrategyMock) Notify(id entities.MachineId, level entities.HealthLevel) error {
	m.wasCalled = true

	if m.shouldReturnError {
		return errReport
	}

	return nil
}

type repositoryMock struct {
	loadedMachine         *entities.Machine
	savedMachine          *entities.Machine
	shouldReturnLoadError bool
	shouldReturnSaveError bool
}

func (r *repositoryMock) Load(id entities.MachineId) (*entities.Machine, error) {
	if r.shouldReturnLoadError {
		return nil, errLoad
	}

	return r.loadedMachine, nil
}

func (r *repositoryMock) Save(machine *entities.Machine) error {
	r.savedMachine = machine

	if r.shouldReturnSaveError {
		return errSave
	}

	return nil
}

const newMachineName string = "the machine"
const existingMachineName string = "the old machine"

var expectedMissingUpdates []entities.MissingUpdate = []entities.MissingUpdate{}
var errReport error = errors.New("report error")
var errLoad error = errors.New("load error")
var errSave error = errors.New("save error")

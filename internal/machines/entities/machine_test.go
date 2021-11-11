package entities

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestReport(t *testing.T) {
	strategyMock := &notificationStrategy{
		reportedHealthLevel: -1,
	}

	var cases = []struct {
		expectedLevel HealthLevel
		severities    []Severity
	}{
		{Healthy, []Severity{Unspecified, Unspecified}},
		{Warning, []Severity{Low, Unspecified}},
		{Warning, []Severity{Low, Low}},
		{Danger, []Severity{Critical, Low}},
		{Danger, []Severity{Critical, Unspecified}},
		{Danger, []Severity{Critical, Critical}},
		{Danger, []Severity{Important, Low}},
		{Danger, []Severity{Important, Unspecified}},
		{Danger, []Severity{Important, Critical}},
		{Danger, []Severity{Important, Important}},
	}

	for _, testCase := range cases {
		expectedId := MachineId(uuid.New())

		machine := Machine{
			h:       health{Healthy},
			missing: []MissingUpdate{},
			Name:    machineName,
			Id:      expectedId,
		}

		expectedMissingUpdates := createMissingUpdates(testCase.severities...)
		err := machine.Report(expectedMissingUpdates, strategyMock)

		if err != nil {
			t.Errorf("Report should not return error %s!", err)
		}

		actualLevel := machine.h.level

		if actualLevel != testCase.expectedLevel {
			t.Errorf("Report should change machine health level to expected %d, but was %d", testCase.expectedLevel, actualLevel)
		}

		if strategyMock.reportedMachineId != expectedId {
			t.Errorf("Strategy was reported with wrong machine id, expected %s, but was %s", expectedId, strategyMock.reportedMachineId)
		}

		if strategyMock.reportedHealthLevel != testCase.expectedLevel {
			t.Errorf("Strategy was reported with wrong health level, expected %d, but was %d", testCase.expectedLevel, strategyMock.reportedHealthLevel)
		}

		if !compareMissingUpdates(machine.missing, expectedMissingUpdates) {
			t.Errorf("Machine should remember missing updates!")
		}
	}
}

func TestReportReturnsError(t *testing.T) {
	strategyMock := &notificationStrategy{
		shouldReturnError: true,
	}

	machine := Machine{
		h:       health{Healthy},
		missing: []MissingUpdate{},
		Name:    machineName,
	}

	expectedMissingUpdates := createMissingUpdates(Important)
	err := machine.Report(expectedMissingUpdates, strategyMock)

	if err == nil {
		t.Errorf("Report should return error!")
	}

	if err != errReport {
		t.Errorf("Error mismatch! Expected %s, but was %s", errReport, err)
	}
}

func TestGetHealthLevel(t *testing.T) {
	machine := Machine{
		h:       health{Warning},
		missing: []MissingUpdate{},
		Name:    machineName,
	}

	actual := machine.GetHealthLevel()

	if actual != Warning {
		t.Errorf("GetMachineHealth level expected as %d, but was %d", Warning, actual)
	}
}

func TestGetMissingUpdates(t *testing.T) {
	expected := createMissingUpdates([]Severity{Critical, Low, Important}...)

	machine := Machine{
		h:       health{Healthy},
		missing: expected,
		Name:    machineName,
	}

	actual := machine.GetMissingUpdates()

	if !compareMissingUpdates(expected, actual) {
		t.Errorf("Machine should remember missing updates!")
	}
}

func TestCreate(t *testing.T) {
	id := MachineId(uuid.New())
	mu := createMissingUpdates([]Severity{Critical, Important}...)
	machine := CreateMachine(id, machineName, mu)

	if machine.Name != machineName {
		t.Errorf("New machine shoud have expected name %s, but was %s", machineName, machine.Name)
	}

	if machine.Id != id {
		t.Errorf("New machine shoud have expected id %s, but was %s", id, machine.Id)
	}

	if !compareMissingUpdates(mu, machine.GetMissingUpdates()) {
		t.Errorf("Machine should remember missing updates!")
	}
}

func BenchmarkReport(b *testing.B) {
	machine := CreateMachine(MachineId(uuid.New()), machineName, []MissingUpdate{})
	mu := createMissingUpdates([]Severity{Critical, Important}...)
	for i := 0; i < b.N; i++ {
		machine.Report(mu, &notificationStrategy{})
	}
}

type notificationStrategy struct {
	reportedMachineId   MachineId
	reportedHealthLevel HealthLevel
	shouldReturnError   bool
}

func (s *notificationStrategy) Notify(machineId MachineId, level HealthLevel) error {
	if s.shouldReturnError {
		return errReport
	}

	s.reportedMachineId = machineId
	s.reportedHealthLevel = level
	return nil
}

func createMissingUpdates(severities ...Severity) (mu []MissingUpdate) {
	for _, s := range severities {
		mu = append(mu, MissingUpdate{
			Severity: s,
		})
	}

	return mu
}

func compareMissingUpdates(first []MissingUpdate, second []MissingUpdate) bool {
	if len(first) != len(second) {
		return false
	}
	for i, v := range first {
		if v != second[i] {
			return false
		}
	}
	return true
}

const machineName string = "the machine"

var errReport error = errors.New("expected error")

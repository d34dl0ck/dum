package entities

import (
	"errors"
	"testing"
)

func TestReport(t *testing.T) {
	strategyMock := &notificationStrategy{
		reportedMachineName: "",
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
		machine := Machine{
			h:       health{Healthy},
			missing: []MissingUpdate{},
			Name:    machineName,
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

		if strategyMock.reportedMachineName != machineName {
			t.Errorf("Strategy was reported with wrong machine name, expected %s, but was %s", machineName, strategyMock.reportedMachineName)
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
	mu := createMissingUpdates([]Severity{Critical, Important}...)
	machine := CreateMachine(machineName, mu)

	if machine.Name != machineName {
		t.Errorf("New machine shoud have expected name %s, but was %s", machineName, machine.Name)
	}

	if !compareMissingUpdates(mu, machine.GetMissingUpdates()) {
		t.Errorf("Machine should remember missing updates!")
	}
}

func BenchmarkReport(b *testing.B) {
	machine := CreateMachine(machineName, []MissingUpdate{})
	mu := createMissingUpdates([]Severity{Critical, Important}...)
	for i := 0; i < b.N; i++ {
		machine.Report(mu, &notificationStrategy{})
	}
}

type notificationStrategy struct {
	reportedMachineName string
	reportedHealthLevel HealthLevel
	shouldReturnError   bool
}

func (s *notificationStrategy) Notify(machineName string, level HealthLevel) error {
	if s.shouldReturnError {
		return errReport
	}

	s.reportedMachineName = machineName
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

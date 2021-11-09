package entities

import "testing"

func TestRecalculate(t *testing.T) {
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
		health := health{
			level: Healthy,
		}

		var missingUpdates []MissingUpdate

		for _, severity := range testCase.severities {
			missingUpdates = append(missingUpdates, MissingUpdate{
				Severity: severity,
			})
		}

		actual := health.Recalculate(missingUpdates)

		if actual.level != testCase.expectedLevel {
			t.Errorf("Health level mismatch! Expected %d but was %d", testCase.expectedLevel, actual.level)
		}
	}
}

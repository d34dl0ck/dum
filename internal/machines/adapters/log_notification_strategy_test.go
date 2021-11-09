package adapters

import (
	"dum/internal/machines/entities"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func TestLogWriting(t *testing.T) {
	file, err := os.Create(testFile)

	if err != nil {
		t.Errorf("Cannot setup test due to error on file creation %s", err)
		return
	}
	defer os.Remove(testFile)
	defer file.Close()

	strategy := NewLogNotificationStrategy(log.New(file, "", 0))
	err = strategy.Notify(machineName, healthLevel)

	if err != nil {
		t.Errorf("Expected nil err, but was %s", err)
		return
	}

	raw, err := os.ReadFile(testFile)

	if err != nil {
		t.Errorf("Cannot read log file %s", err)
		return
	}

	expected := fmt.Sprintf(template, machineName, healthLevel)
	actual := strings.Trim(string(raw), "\t\n\r")

	if actual != expected {
		t.Errorf("Log mismatch! Expected '%s' , but was '%s'", expected, actual)
	}
}

const testFile string = "log_test.txt"
const machineName string = "some machine name"
const healthLevel entities.HealthLevel = entities.Warning

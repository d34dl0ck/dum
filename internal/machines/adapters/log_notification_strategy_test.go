package adapters

import (
	"dum/internal/machines/entities"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestLogWriting(t *testing.T) {
	expectedId := entities.MachineId(uuid.New())
	file, err := os.Create(testFile)

	if err != nil {
		t.Errorf("Cannot setup test due to error on file creation %s", err)
		return
	}
	defer os.Remove(testFile)
	defer file.Close()

	strategy := NewLogNotificationStrategy(log.New(file, "", 0))
	err = strategy.Notify(expectedId, healthLevel)

	if err != nil {
		t.Errorf("Expected nil err, but was %s", err)
		return
	}

	raw, err := os.ReadFile(testFile)

	if err != nil {
		t.Errorf("Cannot read log file %s", err)
		return
	}

	expected := fmt.Sprintf(template, expectedId.String(), healthLevel)
	actual := strings.Trim(string(raw), "\t\n\r")

	if actual != expected {
		t.Errorf("Log mismatch! Expected '%s' , but was '%s'", expected, actual)
	}
}

const testFile string = "log_test.txt"
const healthLevel entities.HealthLevel = entities.Warning

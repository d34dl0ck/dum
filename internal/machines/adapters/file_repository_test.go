package adapters

import (
	"dum/internal/machines/entities"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSuccessSaveLoad(t *testing.T) {
	repo := NewFileRepository()

	file, err := os.Create(RepositoryFileName)

	if err != nil {
		t.Errorf("Cannot setup test due to error %s", err)
		return
	}

	defer os.Remove(RepositoryFileName)
	defer file.Close()
	machine := entities.CreateMachine("testName", []entities.MissingUpdate{
		{
			UpdateId: uuid.New(),
			Severity: entities.Critical,
			Duration: time.Hour,
		},
	})

	err = repo.Save(machine)
	if err != nil {
		t.Errorf("Failed to save machine, because of error %s", err)
		return
	}

	loadedMachine, err := repo.Load(machine.Name)
	if err != nil {
		t.Errorf("Failed to load machine, because of error %s", err)
		return
	}

	if loadedMachine.Name != machine.Name {
		t.Errorf("Machine name mismatch! Expected %s, but was %s", machine.Name, loadedMachine.Name)
	}

	if len(loadedMachine.GetMissingUpdates()) != len(machine.GetMissingUpdates()) {
		t.Errorf("Missing updates length mismatch! Expected %d, but was %d", len(machine.GetMissingUpdates()), len(loadedMachine.GetMissingUpdates()))
	}

	if loadedMachine.GetMissingUpdates()[0].Duration != machine.GetMissingUpdates()[0].Duration {
		t.Errorf("Missing update duration! Expected %s, but was %s", machine.GetMissingUpdates()[0].Duration, loadedMachine.GetMissingUpdates()[0].Duration)
	}

	if loadedMachine.GetMissingUpdates()[0].Severity != machine.GetMissingUpdates()[0].Severity {
		t.Errorf("Missing update severity! Expected %d, but was %d", machine.GetMissingUpdates()[0].Severity, loadedMachine.GetMissingUpdates()[0].Severity)
	}

	if loadedMachine.GetMissingUpdates()[0].UpdateId != machine.GetMissingUpdates()[0].UpdateId {
		t.Errorf("Missing update id! Expected %s, but was %s", machine.GetMissingUpdates()[0].UpdateId, loadedMachine.GetMissingUpdates()[0].UpdateId)
	}

	if loadedMachine.GetHealthLevel() != machine.GetHealthLevel() {
		t.Errorf("Machine health level mismatch! Expected %d, but was %d", machine.GetHealthLevel(), loadedMachine.GetHealthLevel())
	}
}

func TestOptimisticLockError(t *testing.T) {
	repo := NewFileRepository()

	file, err := os.Create(RepositoryFileName)

	if err != nil {
		t.Errorf("Cannot setup test due to error %s", err)
		return
	}

	defer os.Remove(RepositoryFileName)
	defer file.Close()
	machine := entities.CreateMachine("testName", []entities.MissingUpdate{
		{
			UpdateId: uuid.New(),
			Severity: entities.Critical,
			Duration: time.Hour,
		},
	})

	_ = repo.Save(machine)
	repo = NewFileRepository()
	err = repo.Save(machine)

	if err == nil {
		t.Errorf("Optimistic lock doesn't occure!")
		return
	}

	errorBody := err.Error()

	if !strings.HasPrefix(errorBody, "optimistic lock occured!") {
		t.Errorf("Error mismatch! Expected starts with 'optimistic lock occured!', but was something else!")
	}
}

func TestBrokenJsonLoadError(t *testing.T) {
	repo := NewFileRepository()

	file, err := os.Create(RepositoryFileName)

	if err != nil {
		t.Errorf("Cannot setup test due to error %s", err)
		return
	}

	defer os.Remove(RepositoryFileName)
	defer file.Close()
	file.WriteString("not a json")

	machine, err := repo.Load("some name")

	if machine != nil {
		t.Errorf("Machine should be nil!")
	}

	if err == nil {
		t.Error("Error should not be nil!")
	}
}

func TestEmptyLoad(t *testing.T) {
	repo := NewFileRepository()

	file, err := os.Create(RepositoryFileName)

	if err != nil {
		t.Errorf("Cannot setup test due to error %s", err)
		return
	}

	defer os.Remove(RepositoryFileName)
	defer file.Close()
	file.WriteString("[]")

	machine, err := repo.Load("some name")

	if machine != nil {
		t.Errorf("Machine should be nil!")
	}

	if err != nil {
		t.Errorf("Error should be nil, but was %s", err)
	}
}

func TestNoFileLoadError(t *testing.T) {
	repo := NewFileRepository()

	machine, err := repo.Load("some name")

	if machine != nil {
		t.Errorf("Machine should be nil!")
	}

	if err == nil {
		t.Error("Error should not be nil!")
	}
}

func TestNoFileSaveError(t *testing.T) {
	repo := NewFileRepository()

	err := repo.Save(&entities.Machine{})

	if err == nil {
		t.Error("Error should not be nil!")
	}
}

func TestSerializationError(t *testing.T) {
	repo := FileRepository{
		mu:          &sync.Mutex{},
		versionsMap: map[string]MachineVersion{},
		s:           func(i interface{}) ([]byte, error) { return nil, errExpected },
		d:           json.Unmarshal,
		fw:          os.WriteFile,
		fr:          os.ReadFile,
	}

	file, err := os.Create(RepositoryFileName)

	if err != nil {
		t.Errorf("Cannot setup test due to error %s", err)
		return
	}

	defer os.Remove(RepositoryFileName)
	defer file.Close()
	machine := entities.CreateMachine("testName", []entities.MissingUpdate{
		{
			UpdateId: uuid.New(),
			Severity: entities.Critical,
			Duration: time.Hour,
		},
	})

	err = repo.Save(machine)
	if err == nil {
		t.Errorf("Expected error %s, but got no error!", errExpected)
		return
	}

	if err != errExpected {
		t.Errorf("Expected error %s, but was %s", errExpected, err)
	}
}

func TestFileWriteError(t *testing.T) {
	repo := FileRepository{
		mu:          &sync.Mutex{},
		versionsMap: map[string]MachineVersion{},
		s:           json.Marshal,
		d:           json.Unmarshal,
		fw:          func(s string, b []byte, fm os.FileMode) error { return errExpected },
		fr:          os.ReadFile,
	}

	file, err := os.Create(RepositoryFileName)

	if err != nil {
		t.Errorf("Cannot setup test due to error %s", err)
		return
	}

	defer os.Remove(RepositoryFileName)
	defer file.Close()
	machine := entities.CreateMachine("testName", []entities.MissingUpdate{
		{
			UpdateId: uuid.New(),
			Severity: entities.Critical,
			Duration: time.Hour,
		},
	})

	err = repo.Save(machine)
	if err == nil {
		t.Errorf("Expected error %s, but got no error!", errExpected)
		return
	}

	if err != errExpected {
		t.Errorf("Expected error %s, but was %s", errExpected, err)
	}
}

var errExpected error = errors.New("TEST FAIL")

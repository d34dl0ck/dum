package adapters

import (
	"dum/internal/machines/entities"
	"os"
	"testing"
)

func TestShouldCreateFileIfDoesNotExistOnLoad(t *testing.T) {
	defer os.Remove(RepositoryFileName)
	repoMock := repositoryMock{}
	decorator := NewRecoveryFileRepositoryDecorator(&repoMock)

	machine, err := decorator.Load("")

	if machine != expectedMachine {
		t.Errorf("Machine mismatch!")
	}

	if err != nil {
		t.Errorf("Expected error to be nil, but was %s", err)
	}

	raw, err := os.ReadFile(RepositoryFileName)

	if err != nil {
		t.Errorf("Expected error to be nil on result reading, but was %s", err)
	}

	if string(raw) != "[]" {
		t.Error("Unexpected content of result file!")
	}

	if repoMock.isLoaded != true {
		t.Error("Repo wasn't called!")
	}
}

func TestShouldNotCreateFileIfExistsAndNotEmptyOnLoad(t *testing.T) {
	file, err := os.Create(RepositoryFileName)

	if err != nil {
		t.Errorf("Cannot setup test due to error %s", err)
	}

	defer os.Remove(RepositoryFileName)
	file.WriteString("[]")
	file.Close()

	repoMock := repositoryMock{}
	decorator := NewRecoveryFileRepositoryDecorator(&repoMock)

	machine, err := decorator.Load("")

	if machine != expectedMachine {
		t.Errorf("Machine mismatch!")
	}

	if err != nil {
		t.Errorf("Expected error to be nil, but was %s", err)
	}

	if repoMock.isLoaded != true {
		t.Error("Repo wasn't called!")
	}
}

func TestShouldOverwriteFileIfEmptyOnLoad(t *testing.T) {
	file, err := os.Create(RepositoryFileName)

	if err != nil {
		t.Errorf("Cannot setup test due to error %s", err)
	}
	defer os.Remove(RepositoryFileName)
	file.Close()

	repoMock := repositoryMock{}
	decorator := NewRecoveryFileRepositoryDecorator(&repoMock)

	machine, err := decorator.Load("")

	if machine != expectedMachine {
		t.Errorf("Machine mismatch!")
	}

	if err != nil {
		t.Errorf("Expected error to be nil, but was %s", err)
	}

	raw, err := os.ReadFile(RepositoryFileName)

	if err != nil {
		t.Errorf("Expected error to be nil on result reading, but was %s", err)
	}

	if string(raw) != "[]" {
		t.Error("Unexpected content of result file!")
	}

	if repoMock.isLoaded != true {
		t.Error("Repo wasn't called!")
	}
}

func TestShouldCreateFileIfDoesNotExistOnSave(t *testing.T) {
	defer os.Remove(RepositoryFileName)
	repoMock := repositoryMock{}
	decorator := NewRecoveryFileRepositoryDecorator(&repoMock)

	err := decorator.Save(expectedMachine)

	if err != nil {
		t.Errorf("Expected error to be nil, but was %s", err)
	}

	raw, err := os.ReadFile(RepositoryFileName)

	if err != nil {
		t.Errorf("Expected error to be nil on result reading, but was %s", err)
	}

	if string(raw) != "[]" {
		t.Error("Unexpected content of result file!")
	}

	if expectedMachine != repoMock.savedMachine {
		t.Error("Saved machine mismatch!")
	}
}

func TestShouldNotCreateFileIfExistsAndNotEmptyOnSave(t *testing.T) {
	file, err := os.Create(RepositoryFileName)

	if err != nil {
		t.Errorf("Cannot setup test due to error %s", err)
	}

	defer os.Remove(RepositoryFileName)
	file.WriteString("[]")
	file.Close()

	repoMock := repositoryMock{}
	decorator := NewRecoveryFileRepositoryDecorator(&repoMock)

	err = decorator.Save(expectedMachine)

	if err != nil {
		t.Errorf("Expected error to be nil, but was %s", err)
	}

	if expectedMachine != repoMock.savedMachine {
		t.Error("Saved machine mismatch!")
	}
}

func TestShouldOverwriteFileIfEmptyOnSave(t *testing.T) {
	file, err := os.Create(RepositoryFileName)

	if err != nil {
		t.Errorf("Cannot setup test due to error %s", err)
	}
	defer os.Remove(RepositoryFileName)
	file.Close()

	repoMock := repositoryMock{}
	decorator := NewRecoveryFileRepositoryDecorator(&repoMock)

	err = decorator.Save(expectedMachine)

	if err != nil {
		t.Errorf("Expected error to be nil, but was %s", err)
	}

	raw, err := os.ReadFile(RepositoryFileName)

	if err != nil {
		t.Errorf("Expected error to be nil on result reading, but was %s", err)
	}

	if string(raw) == "" {
		t.Error("Unexpected content of result file!")
	}

	if expectedMachine != repoMock.savedMachine {
		t.Error("Saved machine mismatch!")
	}
}

func TestLoadOpenFileError(t *testing.T) {
	decorator := RecoveryFileRepositoryDecorator{
		repo: &repositoryMock{},
		o:    func(s string, i int, fm os.FileMode) (*os.File, error) { return nil, errExpected },
	}
	defer os.Remove(RepositoryFileName)

	_, err := decorator.Load("some_name")

	if err == nil {
		t.Error("Expected err but was nil!")
	}

	if err != errExpected {
		t.Errorf("Expected %s, but was %s", errExpected, err)
	}
}

func TestSaveOpenFileError(t *testing.T) {
	decorator := RecoveryFileRepositoryDecorator{
		repo: &repositoryMock{},
		o:    func(s string, i int, fm os.FileMode) (*os.File, error) { return nil, errExpected },
	}
	defer os.Remove(RepositoryFileName)

	err := decorator.Save(&entities.Machine{})

	if err == nil {
		t.Error("Expected err but was nil!")
	}

	if err != errExpected {
		t.Errorf("Expected %s, but was %s", errExpected, err)
	}
}

func TestFileInfoGetError(t *testing.T) {
	decorator := RecoveryFileRepositoryDecorator{
		repo: &repositoryMock{},
		o:    func(s string, i int, fm os.FileMode) (*os.File, error) { return os.OpenFile(s, i, fm) },
		fir:  func(f *os.File) (os.FileInfo, error) { return nil, errExpected },
	}
	defer os.Remove(RepositoryFileName)

	err := decorator.createFileIfNotExists()

	if err == nil {
		t.Error("Expected err but was nil!")
	}

	if err != errExpected {
		t.Errorf("Expected %s, but was %s", errExpected, err)
	}
}

func TestOsFileWriteError(t *testing.T) {
	decorator := RecoveryFileRepositoryDecorator{
		repo: &repositoryMock{},
		o:    func(s string, i int, fm os.FileMode) (*os.File, error) { return os.OpenFile(s, i, fm) },
		fir:  func(f *os.File) (os.FileInfo, error) { return f.Stat() },
		ofw:  func(f *os.File, s string) (int, error) { return 0, errExpected },
	}
	defer os.Remove(RepositoryFileName)

	err := decorator.createFileIfNotExists()

	if err == nil {
		t.Error("Expected err but was nil!")
	}

	if err != errExpected {
		t.Errorf("Expected %s, but was %s", errExpected, err)
	}
}

type repositoryMock struct {
	savedMachine *entities.Machine
	isLoaded     bool
}

func (r *repositoryMock) Load(name string) (*entities.Machine, error) {
	r.isLoaded = true
	return expectedMachine, nil
}

func (r *repositoryMock) Save(machine *entities.Machine) error {
	r.savedMachine = machine
	return nil
}

var expectedMachine *entities.Machine = &entities.Machine{}

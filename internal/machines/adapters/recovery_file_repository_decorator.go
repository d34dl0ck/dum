package adapters

import (
	"dum/internal/machines/cases"
	"dum/internal/machines/entities"
	"os"
	"sync"
)

// Decorator for file repository, creates file if it doesn't exist
type RecoveryFileRepositoryDecorator struct {
	repo cases.MachineRepository
	mu   *sync.Mutex
	o    func(string, int, os.FileMode) (*os.File, error)
	fir  func(*os.File) (os.FileInfo, error)
	ofw  func(*os.File, string) (int, error)
}

func (r *RecoveryFileRepositoryDecorator) Load(id entities.MachineId) (*entities.Machine, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	err := r.createFileIfNotExists()

	if err != nil {
		return nil, err
	}

	return r.repo.Load(id)
}

func (r *RecoveryFileRepositoryDecorator) Save(machine *entities.Machine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	err := r.createFileIfNotExists()

	if err != nil {
		return err
	}

	return r.repo.Save(machine)
}

func (r *RecoveryFileRepositoryDecorator) createFileIfNotExists() error {
	file, err := r.o(RepositoryFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	info, err := r.fir(file)
	if err != nil {
		return err
	}

	if info.Size() == 0 {
		_, err = r.ofw(file, "{}")

		if err != nil {
			return err
		}
	}

	return nil
}

func NewRecoveryFileRepositoryDecorator(baseRepository cases.MachineRepository) cases.MachineRepository {
	return &RecoveryFileRepositoryDecorator{
		repo: baseRepository,
		o:    os.OpenFile,
		ofw:  func(f *os.File, s string) (int, error) { return f.WriteString(s) },
		fir:  func(f *os.File) (os.FileInfo, error) { return f.Stat() },
		mu:   &sync.Mutex{},
	}
}

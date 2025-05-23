package adapters

import (
	"dum/internal/machines/cases"
	"dum/internal/machines/entities"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

const RepositoryFileName string = "machines.json"

// Function for serialization.
type serializer func(interface{}) ([]byte, error)

// Function for deserialization.
type deserializer func([]byte, interface{}) error

// Function for reading from file
type fileReader func(string) ([]byte, error)

// Function for reading from file
type fileWriter func(string, []byte, os.FileMode) error

// Repository working with file.
type FileRepository struct {
	versionsMap map[string]MachineVersion
	mu          *sync.Mutex
	s           serializer
	d           deserializer
	fr          fileReader
	fw          fileWriter
}

// Machine entity version for changes tracking.
type MachineVersion string

func (r *FileRepository) Load(id entities.MachineId) (*entities.Machine, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	dtoSet, err := r.loadAll()
	if err != nil {
		return nil, err
	}

	for _, dto := range dtoSet {
		if dto.Id == id.String() {
			r.versionsMap[dto.Id] = MachineVersion(dto.Version)
			return dto.toMachine(), nil
		}
	}

	return nil, nil
}

func (r *FileRepository) Save(machine *entities.Machine) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	dtoSet, err := r.loadAll()
	if err != nil {
		return err
	}

	for _, dto := range dtoSet {
		if dto.Id == machine.Id.String() && MachineVersion(dto.Version) != r.versionsMap[machine.Id.String()] {
			return fmt.Errorf("optimistic lock occured! Expected %s version, but was %s version", r.versionsMap[machine.Name], dto.Version)
		}
	}

	missingUpdateDtoSet := []missingUpdateDto{}
	for _, update := range machine.GetMissingUpdates() {
		dto := missingUpdateDto{
			UpdateId: update.UpdateId.String(),
			Duration: update.Duration,
			Severity: int(update.Severity),
		}
		missingUpdateDtoSet = append(missingUpdateDtoSet, dto)
	}
	machineDto := machineDto{
		Name:           machine.Name,
		MissingUpdates: missingUpdateDtoSet,
		Version:        uuid.NewString(),
		Id:             machine.Id.String(),
	}
	dtoSet[machine.Id.String()] = machineDto

	raw, err := r.s(&dtoSet)
	if err != nil {
		return err
	}

	err = r.fw(RepositoryFileName, raw, os.ModeAppend)
	if err != nil {
		return err
	}
	r.versionsMap[machine.Id.String()] = MachineVersion(machineDto.Version)
	return nil
}

func (r *FileRepository) loadAll() (map[string]machineDto, error) {
	raw, err := r.fr(RepositoryFileName)
	if err != nil {
		return nil, err
	}

	if len(raw) == 0 {
		return map[string]machineDto{}, nil
	}

	var dtoSet map[string]machineDto
	err = r.d(raw, &dtoSet)

	if err != nil {
		return nil, err
	}

	return dtoSet, nil
}

func NewFileRepository() cases.MachineRepository {
	return &FileRepository{
		versionsMap: map[string]MachineVersion{},
		mu:          &sync.Mutex{},
		s:           json.Marshal,
		d:           json.Unmarshal,
		fr:          os.ReadFile,
		fw:          os.WriteFile,
	}
}

// Data transfer object for missing update entity.
type missingUpdateDto struct {
	UpdateId string
	Severity int
	Duration time.Duration
}

func (m missingUpdateDto) toMissingUpdate() entities.MissingUpdate {
	return entities.MissingUpdate{
		UpdateId: uuid.MustParse(m.UpdateId),
		Duration: m.Duration,
		Severity: entities.Severity(m.Severity),
	}
}

// Data transfer object for machine entity.
type machineDto struct {
	Id, Name, Version string
	MissingUpdates    []missingUpdateDto
}

func (m machineDto) toMachine() *entities.Machine {
	var missingUpdates []entities.MissingUpdate

	for _, dto := range m.MissingUpdates {
		missingUpdates = append(missingUpdates, dto.toMissingUpdate())
	}

	return entities.CreateMachine(
		entities.MachineId(uuid.MustParse(m.Id)),
		m.Name,
		missingUpdates,
	)
}

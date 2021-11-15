package ports

import (
	"dum/internal/machines/cases"
	"dum/internal/machines/entities"
	"dum/pkg/machines/report"
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type ReportHandler struct {
	strategy           entities.HealthNotificationStrategy
	repo               cases.MachineRepository
	machineNamePattern regexp.Regexp
	commandChan        chan<- cases.Command
}

func (h *ReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.reportMissingUpdates(w, r)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (h *ReportHandler) reportMissingUpdates(w http.ResponseWriter, r *http.Request) {
	matches := h.machineNamePattern.FindStringSubmatch(r.URL.Path)
	if len(matches) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id, err := uuid.Parse(matches[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(r.Body)
	var request report.ReportRequest
	err = dec.Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var missingUpdates []entities.MissingUpdate

	for _, dto := range request.MissingUpdates {
		missingUpdate, err := convert(dto)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		missingUpdates = append(missingUpdates, missingUpdate)
	}

	command := cases.ReportCommand{
		MachineName:          request.MachineName,
		Repository:           h.repo,
		NotificationStrategy: h.strategy,
		MissingUpdates:       missingUpdates,
		MachineId:            entities.MachineId(id),
	}

	h.commandChan <- &command
	w.WriteHeader(http.StatusAccepted)
}

func NewReportHandler(s entities.HealthNotificationStrategy, r cases.MachineRepository, c chan<- cases.Command) *ReportHandler {
	return &ReportHandler{
		strategy:           s,
		repo:               r,
		machineNamePattern: *regexp.MustCompile(`^/api/v1/machines/(.*)/report`),
		commandChan:        c,
	}
}

func convert(dto report.MissingUpdate) (entities.MissingUpdate, error) {
	duration, err := time.ParseDuration(dto.Duration)

	if err != nil {
		return entities.MissingUpdate{}, err
	}

	return entities.MissingUpdate{
		UpdateId: uuid.MustParse(dto.UpdateId),
		Duration: duration,
		Severity: entities.Severity(dto.Severity),
	}, nil
}

package main

import (
	"dum/internal/machines/cases"
	"dum/internal/machines/ports"
	"testing"
)

func TestReturnHandler(t *testing.T) {
	handler := createHandler(make(chan cases.Command))

	if _, ok := handler.(*ports.ReportHandler); !ok {
		t.Errorf("Handler type mismatch!")
	}
}

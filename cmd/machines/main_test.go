package main

import (
	"dum/internal/machines/ports"
	"testing"
)

func TestReturnHandler(t *testing.T) {
	handler := createHandler()

	if _, ok := handler.(*ports.ReportHandler); !ok {
		t.Errorf("Handler type mismatch!")
	}
}

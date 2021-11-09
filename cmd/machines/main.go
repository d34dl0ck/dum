package main

import (
	"dum/internal/machines/adapters"
	"dum/internal/machines/ports"
	"log"
	"net/http"
)

func main() {
	http.Handle("/api/v1/machines/", createHandler())
	log.Default().Println("Ready to listen!")
	http.ListenAndServe(":3000", nil)
}

func createHandler() http.Handler {
	s := adapters.NewLogNotificationStrategy(log.Default())
	r := adapters.NewFileRepository()
	d := adapters.NewRecoveryFileRepositoryDecorator(r)
	return ports.NewReportHandler(s, d)
}

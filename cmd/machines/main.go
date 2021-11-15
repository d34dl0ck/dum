package main

import (
	"context"
	"dum/internal/machines/adapters"
	"dum/internal/machines/cases"
	"dum/internal/machines/ports"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	commandChan := make(chan cases.Command, 100)
	processingGroup := &sync.WaitGroup{}

	processingCtx, cancelProcessing := context.WithCancel(context.Background())
	startProcessing(processingCtx, commandChan, processingGroup)

	httpServer := createServer(commandChan)
	startServer(httpServer)

	waitForOsSignal()

	gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(gracefullCtx); err != nil {
		log.Default().Printf("Shutdown error: %s", err)
		os.Exit(1)
	}

	log.Default().Println("Cancelling processing context...")
	cancelProcessing()
	log.Default().Println("Closing command channel ...")
	close(commandChan)
	log.Default().Println("Waiting for processing group ...")
	processingGroup.Wait()
	log.Default().Printf("Service is gracefully stopped!")

	os.Exit(0)
}

func createServer(c chan cases.Command) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/api/v1/machines/", createHandler(c))
	httpServer := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	return httpServer
}

func startProcessing(ctx context.Context, c chan cases.Command, wg *sync.WaitGroup) {
	for i := 0; i < 4; i++ {
		p := cases.NewReportProcessor(c, log.Default())
		p.Start(ctx, wg)
	}
}

func startServer(s *http.Server) {
	go func() {
		log.Default().Println("Ready to listen!")
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Default().Fatalf("HTTP server ListenAndServe failed unexpectedly: %s", err)
		}
	}()
}

func createHandler(c chan cases.Command) http.Handler {
	// s := adapters.NewLogNotificationStrategy(log.Default())
	s := adapters.RabbitNotifictionStrategy{
		RabbitUrl: "amqp://guest:guest@172.17.0.2:5672/",
	}
	r := adapters.NewFileRepository()
	d := adapters.NewRecoveryFileRepositoryDecorator(r)

	return ports.NewReportHandler(s, d, c)
}

func waitForOsSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	s := <-signalChan
	log.Default().Printf("Got signal %s - shutting down...", s)

	go func() {
		s := <-signalChan
		log.Default().Fatalf("Got signal %s - terminating...", s)
	}()
}

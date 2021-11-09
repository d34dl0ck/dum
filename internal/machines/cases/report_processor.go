package cases

import (
	"context"
	"log"
	"sync"
)

// Processing commands from channel.
type ReportProcessor struct {
	commandChan <-chan Command
	logger      *log.Logger
}

func NewReportProcessor(commandChannel <-chan Command, logger *log.Logger) *ReportProcessor {
	return &ReportProcessor{
		commandChan: commandChannel,
		logger:      logger,
	}
}

// Starts processing from previously selected channel.
func (p *ReportProcessor) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				p.logger.Println("Stopping processor...")
				for command := range p.commandChan {
					p.execute(command)
				}
				p.logger.Println("Stopped processor!")
				return
			case command, ok := <-p.commandChan:
				if ok {
					p.execute(command)
				}
			}
		}
	}()
}

func (r *ReportProcessor) execute(c Command) {
	err := c.Execute()

	if err != nil {
		r.logger.Printf("Got an error while executing command - %s", err)
	} else {
		r.logger.Println("Successfuly executed command!")
	}
}

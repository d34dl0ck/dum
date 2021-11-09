package cases

import (
	"context"
	"errors"
	"log"
	"sync"
	"testing"
	"time"
)

func TestReportProcessorExecutesCommands(t *testing.T) {
	const commandCount = 5
	commandChannel := make(chan Command, 4)
	ctx, cancel := context.WithCancel(context.Background())
	p := NewReportProcessor(commandChannel, log.Default())
	wg := &sync.WaitGroup{}
	mocks := []*commandMock{}

	for i := 0; i < commandCount; i++ {
		mocks = append(mocks, &commandMock{false, i%2 == 1})
	}

	p.Start(ctx, wg)

	for _, c := range mocks {
		commandChannel <- c
	}

	close(commandChannel)
	cancel()
	wg.Wait()

	for _, m := range mocks {
		if !m.wasExecuted {
			t.Errorf("Command should be executed, but it wasn't!")
		}
	}
}

type commandMock struct {
	wasExecuted, shouldReturnError bool
}

func (c *commandMock) Execute() error {
	c.wasExecuted = true

	if c.shouldReturnError {
		time.Sleep(time.Millisecond * 200)
		return errors.New("some error")
	}

	return nil
}

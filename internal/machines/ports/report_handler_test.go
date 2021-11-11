package ports

import (
	"dum/internal/machines/cases"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestNotFoundIfCannotFindMachineNameInUrl(t *testing.T) {
	handler := NewReportHandler(nil, nil, make(chan<- cases.Command))
	writerMock := responseWriter{
		c: &writerResultContainer{},
	}

	url, _ := url.Parse("/api/v1/machindawdadwa")
	handler.ServeHTTP(writerMock, &http.Request{
		URL:    url,
		Method: http.MethodPost,
	})

	if writerMock.c.writtenStatusCode != 404 {
		t.Errorf("Response code mismatch! Expected %d, but was %d!", 404, writerMock.c.writtenStatusCode)
	}
}

func TestBadRequestIfRequestHasNotValidId(t *testing.T) {
	handler := NewReportHandler(nil, nil, make(chan<- cases.Command))
	writerMock := responseWriter{
		c: &writerResultContainer{},
	}

	url, _ := url.Parse("/api/v1/machines/test/report")
	handler.ServeHTTP(writerMock, &http.Request{
		URL:    url,
		Method: http.MethodPost,
		Body:   io.NopCloser(strings.NewReader("not_a_json")),
	})

	if writerMock.c.writtenStatusCode != 400 {
		t.Errorf("Response code mismatch! Expected %d, but was %d!", 400, writerMock.c.writtenStatusCode)
	}
}

func TestBadRequestIfRequestHasNotValidJsonBody(t *testing.T) {
	handler := NewReportHandler(nil, nil, make(chan<- cases.Command))
	writerMock := responseWriter{
		c: &writerResultContainer{},
	}

	url, _ := url.Parse("/api/v1/machines/test/report")
	handler.ServeHTTP(writerMock, &http.Request{
		URL:    url,
		Method: http.MethodPost,
		Body:   io.NopCloser(strings.NewReader("not_a_json")),
	})

	if writerMock.c.writtenStatusCode != 400 {
		t.Errorf("Response code mismatch! Expected %d, but was %d!", 400, writerMock.c.writtenStatusCode)
	}
}

func TestBadRequestIfCannotDeserializeDto(t *testing.T) {
	handler := NewReportHandler(nil, nil, make(chan<- cases.Command))
	writerMock := responseWriter{
		c: &writerResultContainer{},
	}

	url, _ := url.Parse("/api/v1/machines/1a3fccff-2d7b-45f0-a3c4-50a7bb50d06d/report")
	handler.ServeHTTP(writerMock, &http.Request{
		URL:    url,
		Method: http.MethodPost,
		Body:   io.NopCloser(strings.NewReader("{ \"MachineName\": \"test\", \"MissingUpdates\": [{ \"duration\": \"dadwwadwa\", \"updateId\": \"1a3fccff-2d7b-45f0-a3c4-50a7bb50d06c\", \"severity\": 2 }] }")),
	})

	if writerMock.c.writtenStatusCode != 400 {
		t.Errorf("Response code mismatch! Expected %d, but was %d!", 400, writerMock.c.writtenStatusCode)
	}
}

func TestAccepted(t *testing.T) {
	c := make(chan cases.Command, 1)
	handler := NewReportHandler(nil, nil, c)
	writerMock := responseWriter{
		c: &writerResultContainer{},
	}

	url, _ := url.Parse("/api/v1/machines/1a3fccff-2d7b-45f0-a3c4-50a7bb50d06e/report")
	handler.ServeHTTP(writerMock, &http.Request{
		URL:    url,
		Method: http.MethodPost,
		Body:   io.NopCloser(strings.NewReader("{ \"MachineName\": \"test\", \"MissingUpdates\": [{ \"duration\": \"30s\", \"updateId\": \"1a3fccff-2d7b-45f0-a3c4-50a7bb50d06c\", \"severity\": 2 }] }")),
	})

	if writerMock.c.writtenStatusCode != 202 {
		t.Errorf("Response code mismatch! Expected %d, but was %d!", 202, writerMock.c.writtenStatusCode)
	}

	command, isOpen := <-c

	if !isOpen {
		t.Errorf("Channel should not be closed!")
	}

	if command == nil {
		t.Errorf("Command should not be nil!")
	}
}

func TestNotImplementedIfNotPostMethod(t *testing.T) {
	handler := NewReportHandler(nil, nil, make(chan<- cases.Command))
	writerMock := responseWriter{
		c: &writerResultContainer{},
	}

	methods := []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPut,
		http.MethodTrace,
	}

	for _, method := range methods {
		writerMock.c.writtenStatusCode = 0
		handler.ServeHTTP(writerMock, &http.Request{Method: method})

		if writerMock.c.writtenStatusCode != 501 {
			t.Errorf("Response code mismatch! Expected %d, but was %d!", 501, writerMock.c.writtenStatusCode)
		}
	}
}

func BenchmarkHandler(b *testing.B) {
	handler := NewReportHandler(nil, nil, make(chan<- cases.Command, b.N))
	writerMock := responseWriter{
		c: &writerResultContainer{},
	}

	url, _ := url.Parse("/api/v1/machines/test/report")

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(writerMock, &http.Request{
			URL:    url,
			Method: http.MethodPost,
			Body:   io.NopCloser(strings.NewReader("[{ \"duration\": \"30s\", \"updateId\": \"1a3fccff-2d7b-45f0-a3c4-50a7bb50d06c\", \"severity\": 2 }]")),
		})
	}
}

type writerResultContainer struct {
	writtenStatusCode int
}

type responseWriter struct {
	c *writerResultContainer
}

func (r responseWriter) WriteHeader(statusCode int) {
	r.c.writtenStatusCode = statusCode
}

func (r responseWriter) Write([]byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (r responseWriter) Header() http.Header {
	return nil
}

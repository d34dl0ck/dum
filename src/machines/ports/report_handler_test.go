package ports

import (
	"dum/machines/entities"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestNotFoundIfCannotFindMachineNameInUrl(t *testing.T) {
	handler := NewReportHandler(nil, nil)
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

func TestBadRequestIfRequestHasNotValidJsonBody(t *testing.T) {
	handler := NewReportHandler(nil, nil)
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
	handler := NewReportHandler(nil, nil)
	writerMock := responseWriter{
		c: &writerResultContainer{},
	}

	url, _ := url.Parse("/api/v1/machines/test/report")
	handler.ServeHTTP(writerMock, &http.Request{
		URL:    url,
		Method: http.MethodPost,
		Body:   io.NopCloser(strings.NewReader("[{ \"duration\": \"dadwwadwa\", \"updateId\": \"1a3fccff-2d7b-45f0-a3c4-50a7bb50d06c\", \"severity\": 2 }]")),
	})

	if writerMock.c.writtenStatusCode != 400 {
		t.Errorf("Response code mismatch! Expected %d, but was %d!", 400, writerMock.c.writtenStatusCode)
	}
}

func TestAccepted(t *testing.T) {
	handler := NewReportHandler(notificationStrategy{}, repository{})
	writerMock := responseWriter{
		c: &writerResultContainer{},
	}

	url, _ := url.Parse("/api/v1/machines/test/report")
	handler.ServeHTTP(writerMock, &http.Request{
		URL:    url,
		Method: http.MethodPost,
		Body:   io.NopCloser(strings.NewReader("[{ \"duration\": \"30s\", \"updateId\": \"1a3fccff-2d7b-45f0-a3c4-50a7bb50d06c\", \"severity\": 2 }]")),
	})

	if writerMock.c.writtenStatusCode != 202 {
		t.Errorf("Response code mismatch! Expected %d, but was %d!", 202, writerMock.c.writtenStatusCode)
	}
}

func TestInternalServerError(t *testing.T) {
	handler := NewReportHandler(notificationStrategy{}, repository{true})
	writerMock := responseWriter{
		c: &writerResultContainer{},
	}

	url, _ := url.Parse("/api/v1/machines/test/report")
	handler.ServeHTTP(writerMock, &http.Request{
		URL:    url,
		Method: http.MethodPost,
		Body:   io.NopCloser(strings.NewReader("[{ \"duration\": \"30s\", \"updateId\": \"1a3fccff-2d7b-45f0-a3c4-50a7bb50d06c\", \"severity\": 2 }]")),
	})

	if writerMock.c.writtenStatusCode != 500 {
		t.Errorf("Response code mismatch! Expected %d, but was %d!", 500, writerMock.c.writtenStatusCode)
	}
}

func TestNotImplementedIfNotPostMethod(t *testing.T) {
	handler := NewReportHandler(nil, nil)
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

type repository struct {
	shouldFail bool
}

func (r repository) Load(name string) (*entities.Machine, error) {
	if r.shouldFail {
		return nil, errors.New("some error")
	}

	return nil, nil
}

func (r repository) Save(machine *entities.Machine) error {
	return nil
}

type notificationStrategy struct{}

func (s notificationStrategy) Notify(machineName string, level entities.HealthLevel) error {
	return nil
}

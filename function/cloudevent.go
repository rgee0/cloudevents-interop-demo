package function

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/docker/distribution/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/openfaas-incubator/go-function-sdk"
)

// CloudEvent v0.1
// https://github.com/cloudevents/spec/blob/v0.1/json-format.md
type CloudEvent struct {
	Type             string            `json:"type"`
	EventTypeVersion string            `json:"eventTypeVersion,omitempty"`
	SpecVersion      string            `json:"specVersion"`
	Source           string            `json:"source"`
	ID               string            `json:"id"`
	Time             time.Time         `json:"time,omitempty"`
	RelatedID        string            `json:"relatedid,omitempty"`
	ContentType      string            `json:"contentType,omitempty"`
	Extensions       map[string]string `json:"extensions,omitempty"`
	Data             json.RawMessage   `json:"data,omitempty"`
}

const headerPrefix = "ce-"

func (c *CloudEvent) inititialise(eType string) {

	c.Type = eType
	c.SpecVersion = "0.1"
	c.Source = "https://rgee0.o6s.io/cloudevents-interop-demo"
	c.ID = uuid.Generate().String()
	c.Time = time.Now()
	c.ContentType = "application/json"
}

func getCloudEvent(req *handler.Request, structuredRequest bool) (*CloudEvent, error) {

	if structuredRequest {
		return getStructuredCloudEvent(req.Body)
	}

	return getBinaryCloudEvent(req.Header)
}

// getStructuredCloudEvent returns a pointer to a CloudEvent extracted from the
// structured request submitted to the handler
func getStructuredCloudEvent(req []byte) (*CloudEvent, error) {
	c := CloudEvent{}

	if err := json.Unmarshal(req, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

// getBinaryCloudEvent returns a pointer to a CloudEvent extracted from the
// binary request submitted to the handler
func getBinaryCloudEvent(header map[string][]string) (*CloudEvent, error) {
	c := CloudEvent{}

	var headers = make(map[string]string)

	for headerKey, headerVal := range header {

		if !strings.HasPrefix(headerKey, headerPrefix) {
			continue
		}

		headerKey = strings.TrimPrefix(headerKey, headerPrefix)
		headerKey = strings.Replace(headerKey, "-", "", -1)
		headers[headerKey] = headerVal[0]

	}

	mapstructure.Decode(headers, &c)

	return &c, nil
}

func setStructureCloudEvent(c *CloudEvent) ([]byte, error) {
	retBytes, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return retBytes, nil
}

func cloudEventResponse(c *CloudEvent, structuredRequest bool) (handler.Response, error) {

	var (
		err      error
		bMessage []byte
	)

	if structuredRequest {

		bMessage, err = setStructureCloudEvent(c)

		return handler.Response{
			Body:       bMessage,
			StatusCode: http.StatusOK,
			Header: map[string][]string{
				"Content-Type": []string{"application/cloudevents+json"},
			},
		}, err
	}

	return handler.Response{
		Body:       bMessage,
		StatusCode: http.StatusOK,
	}, err
}

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

func initCloudEvent(eType string, data map[string]string, reqID string) *CloudEvent {

	dataField, err := json.Marshal(&data)

	if err != nil {
		dataField = nil
	}

	return &CloudEvent{
		Type:        eType,
		SpecVersion: "0.1",
		Source:      "https://rgee0.o6s.io/cloudevents-interop-demo",
		ID:          uuid.Generate().String(),
		RelatedID:   reqID,
		Time:        time.Now(),
		ContentType: "application/json",
		Data:        dataField,
	}
}

func getCloudEvent(req *handler.Request, structuredRequest bool) (*CloudEvent, error) {

	if structuredRequest {
		return getStructuredCloudEvent(req.Body)
	}

	c, err := getBinaryCloudEvent(req.Header)
	c.Data = req.Body
	return c, err
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

		if !strings.EqualFold(headerKey[:3], headerPrefix) {
			continue
		}

		headerKey = headerKey[3:]
		headerKey = strings.Replace(headerKey, "-", "", -1)
		headers[headerKey] = headerVal[0]

	}

	mapstructure.Decode(headers, &c)

	return &c, nil
}

func setCloudEvent(c *CloudEvent) ([]byte, error) {
	retBytes, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return retBytes, nil
}

func sendCloudEvent(c *CloudEvent, err error) (handler.Response, error) {

	if err != nil {
		return handler.Response{}, err
	}

	bMessage, err := setCloudEvent(c)

	return handler.Response{
		Body:       bMessage,
		StatusCode: http.StatusOK,
		Header: map[string][]string{
			"Content-Type": []string{"application/cloudevents+json"},
		},
	}, err
}

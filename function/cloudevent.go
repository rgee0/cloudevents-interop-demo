package function

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/docker/distribution/uuid"
	"github.com/mitchellh/mapstructure"
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

func initCloudEvent(eType string) CloudEvent {
	return CloudEvent{
		Type:        eType,
		SpecVersion: "0.1",
		Source:      "https://rgee0.o6s.io/cloudevents-interop-demo",
		ID:          uuid.Generate().String(),
		Time:        time.Now(),
		ContentType: "application/json",
	}
}

// getCloudEvent returns a pointer to a CloudEvent extracted from the
// request submitted to the handler
func getCloudEvent(req []byte) (*CloudEvent, error) {
	c := CloudEvent{}

	if err := json.Unmarshal(req, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func getBinaryCloudEvent(header map[string][]string) (*CloudEvent, error) {
	c := CloudEvent{}

	var headers = make(map[string]string)

	for headerKey, headerVal := range header {

		if !strings.HasPrefix(headerKey, "ce-") {
			continue
		}

		headerKey = strings.TrimPrefix(headerKey, "ce-")
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

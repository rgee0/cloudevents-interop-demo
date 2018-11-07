package function

import (
	"encoding/json"
	"time"

	"github.com/docker/distribution/uuid"
)

// CloudEvent v0.1
// https://github.com/cloudevents/spec/blob/v0.1/json-format.md
type CloudEvent struct {
	EventType          string            `json:"eventType"`
	EventTypeVersion   string            `json:"eventTypeVersion,omitempty"`
	CloudEventsVersion string            `json:"cloudEventsVersion"`
	Source             string            `json:"source"`
	EventID            string            `json:"eventID"`
	EventTime          time.Time         `json:"eventTime,omitempty"`
	ContentType        string            `json:"contentType,omitempty"`
	Extensions         map[string]string `json:"extensions,omitempty"`
	Data               json.RawMessage   `json:"data,omitempty"`
}

func initCloudEvent(eventType string) CloudEvent {
	return CloudEvent{
		EventType:          eventType,
		CloudEventsVersion: "0.1",
		Source:             "https://rgee0.o6s.io/cloudevents-interop-demo",
		EventID:            uuid.Generate().String(),
		EventTime:          time.Now(),
		ContentType:        "application/json",
	}
}

func getCloudEvent(req []byte) (*CloudEvent, error) {
	c := CloudEvent{}
	if err := json.Unmarshal(req, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func setCloudEvent(c *CloudEvent) ([]byte, error) {
	retBytes, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return retBytes, nil
}

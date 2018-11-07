package function

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/openfaas-incubator/go-function-sdk"
)

const (
	structuredContentType = "cloudevents"
	wordsURLEnvVar        = "wordsURL"
)

var wordList = make(map[string][]string)

func init() {

	wordList = getWordList()

}

// The receiver of the event can distinguish between the two modes by inspecting the Content-Type header value.
// If the value is prefixed with the CloudEvents media type application/cloudevents, indicating the use of a known
// event format, the receiver uses structured mode, otherwise it defaults to binary mode.
// https://github.com/cloudevents/spec/blob/a12b6b618916c89bfa5595fc76732f07f89219b5/http-transport-binding.md#3-http-message-mapping
func isStructured(httpContentTypes []string) bool {

	for _, cType := range httpContentTypes {
		if strings.Contains(cType, structuredContentType) {
			return true
		}
	}
	return false
}

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {

	var (
		err      error
		bMessage []byte
		c        *CloudEvent
		debug    string
	)

	if len(wordList) == 0 {
		wordList = getWordList()
	}

	if isStructured(req.Header["Content-Type"]) {

		c, err = getCloudEvent(req.Body)
		debug = "structured"

	} else {

		c, err = getBinaryCloudEvent(req.Header)
		c.Data = req.Body
		debug = "binary"

	}

	wordType := strings.Split(c.EventType, ".")[2]
	dataField := getWordValue(wordList[wordType])

	if dataField != nil {
		retEventType := strings.Replace(c.EventType, "found", "picked", -1)
		retEvent := initCloudEvent(retEventType)
		retEvent.Data, err = json.Marshal(&dataField)
		bMessage, err = setCloudEvent(&retEvent)
	}

	return handler.Response{
		Body:       bMessage,
		StatusCode: http.StatusOK,
		Header: map[string][]string{
			"Content-Type": []string{"application/cloudevents+json"},
			"Debug":        []string{debug},
		},
	}, err
}

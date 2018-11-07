package function

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/openfaas-incubator/go-function-sdk"
)

var wordList = make(map[string][]string)

func init() {

	wordList = getWordList()

}

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {

	var err error
	var bMessage []byte

	if len(wordList) == 0 {
		wordList = getWordList()
	}

	c, err := getCloudEvent(req.Body)

	wordType := strings.Split(c.EventType, ".")[2]
	dataField := getWordValue(wordList[wordType])

	if dataField != nil {
		returnEventType := strings.Replace(c.EventType, "found", "picked", -1)
		retCEvent := initCloudEvent(returnEventType)
		retCEvent.Data, err = json.Marshal(&dataField)
		bMessage, err = setCloudEvent(&retCEvent)
	}

	return handler.Response{
		Body:       bMessage,
		StatusCode: http.StatusOK,
		Header: map[string][]string{
			"Content-Type": []string{"application/cloudevents+json"},
		},
	}, err
}

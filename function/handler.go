package function

import (
	"math/rand"
	"strings"
	"time"

	"github.com/openfaas-incubator/go-function-sdk"
)

const (
	structuredContentType = "cloudevents"
	wordsURLEnvVar        = "wordsURL"
	reqEventTypePattern   = "found"
	resEventTypePattern   = "picked"
)

var wordList = make(map[string][]string)

func init() {

	wordList = getWordList()
	rand.Seed(time.Now().UTC().UnixNano())

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
		err         error
		c, retEvent *CloudEvent
	)

	if len(wordList) == 0 {
		wordList = getWordList()
	}

	structuredRequest := isStructured(req.Header["Content-Type"])
	c, err = getCloudEvent(&req, structuredRequest)

	//temporary
	/*if !structuredRequest {
		postBackBody, _ := json.Marshal(&req)
		postBack, _ := http.NewRequest("POST", "http://requestbin.fullcontact.com/1ijmli01", bytes.NewBuffer(postBackBody))
		client := &http.Client{}
		postBack.Header.Set("Content-Type", "application/json")
		_, _ = client.Do(postBack)
	}*/
	//temporary

	wordType := strings.Split(c.Type, ".")[2]
	dataVal := getWordValue(wordList[wordType])

	if dataVal != nil {
		retEventType := strings.Replace(c.Type, reqEventTypePattern, resEventTypePattern, -1)
		retEvent = initCloudEvent(retEventType, dataVal, c.ID)
	}

	return sendCloudEvent(retEvent, structuredRequest, true, err)

}

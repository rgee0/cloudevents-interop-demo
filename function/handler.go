package function

import (
	"math/rand"
	"net/http"
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

// sendCloudEvent - take an existing cloud event struct and generate the handler response for it according to
// the demo conventions.  Respond to requests with the respective event type (binary/structured).
// If X-Callback-URL is set then send only a 202 to the client with the response event sent to X-Callback-URL
func sendCloudEvent(c *CloudEvent, structuredRequest bool, callbackURL string, err error) (handler.Response, error) {

	var (
		bMessage   []byte
		headerVals map[string][]string
		statusCode = http.StatusOK
	)

	if err != nil {
		return handler.Response{}, err
	}

	if structuredRequest {
		bMessage, headerVals, err = setStructuredCloudEvent(c)
	} else {
		bMessage, headerVals, err = setBinaryCloudEvent(c)
	}

	if len(callbackURL) > 0 {
		statusCode = http.StatusAccepted
	}

	return handler.Response{
		Body:       bMessage,
		StatusCode: statusCode,
		Header:     headerVals,
	}, err
}

//temporary
/*if !structuredRequest {
	postBackBody, _ := json.Marshal(&req)
	postBack, _ := http.NewRequest("POST", "http://requestbin.fullcontact.com/1ijmli01", bytes.NewBuffer(postBackBody))
	client := &http.Client{}
	postBack.Header.Set("Content-Type", "application/json")
	_, _ = client.Do(postBack)
}*/
//temporary

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {

	var (
		err         error
		c, retEvent *CloudEvent
		callbackURL string
	)

	if len(wordList) == 0 {
		wordList = getWordList()
	}

	structuredRequest := isStructured(req.Header["Content-Type"])

	if cbVal, ok := req.Header["X-Callback-URL"]; ok {
		callbackURL = cbVal[0]
	}

	c, err = getCloudEvent(&req, structuredRequest)

	wordType := strings.Split(c.Type, ".")[2]
	dataVal := getWordValue(wordList[wordType])

	if dataVal != nil {
		retEventType := strings.Replace(c.Type, reqEventTypePattern, resEventTypePattern, -1)
		retEvent = initCloudEvent(retEventType, dataVal, c.ID)
	}

	return sendCloudEvent(retEvent, structuredRequest, callbackURL, err)

}

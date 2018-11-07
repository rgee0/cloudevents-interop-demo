package function

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func getWordList() map[string][]string {

	var wordMap map[string][]string

	var netClient = &http.Client{
		Timeout: time.Second * 3,
	}

	wordsURL := os.Getenv(wordsURLEnvVar)

	if len(wordsURL) <= 0 {
		log.Panic("wordsURL env var not set or empty")
	}

	resp, getErr := netClient.Get(wordsURL)

	if getErr != nil {
		panic(getErr.Error())
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		panic(readErr.Error())
	}

	parseErr := json.Unmarshal(body, &wordMap)
	if parseErr != nil {
		panic(parseErr.Error())
	}

	return wordMap
}

func getWordValue(wordList []string) map[string]string {

	if listLength := len(wordList) - 1; listLength > 0 {
		rand.Seed(time.Now().Unix())
		return map[string]string{"word": wordList[rand.Intn(listLength)]}
	}
	return nil
}

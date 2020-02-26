package aftg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Connector struct {
	ApiKey string
	RetryAmount int8
}

func (aftg *Connector) runAftgRequest(method string, path string, requestBody io.Reader, queryParams map[string]string, additionalHeaders map[string]string) (int, []byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, os.Getenv("AFTG_API_URL") + path, requestBody)
	if err != nil {
		return -1, nil, err
	}
	req.Header.Add("X-API-KEY", aftg.ApiKey)
	req.Header.Add("Content-Type", "application/json")
	for key, value := range additionalHeaders {
		req.Header.Add(key, value)
	}

	query := req.URL.Query()
	for key, value := range queryParams {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, nil, err
	}
	return resp.StatusCode, body, err
}

type NTP struct {
	SrvReceptionTime int64 `json:"srvReceptionTime"`
	ClientTransmissionTime int64 `json:"clientTransmissionTime"`
	SrvTransmissionTime int64 `json:"srvTransmissionTime"`
	ClientReceptionTime int64 `json:"clientReceptionTime"`
}

func (aftg *Connector) GetSrvDelay() int64 {
	var ntp NTP

	code, body, err := aftg.runAftgRequest("GET", "ntp", nil,
		map[string]string{"clientTransmissionTime": strconv.FormatInt(time.Now().UnixNano() / int64(time.Millisecond), 10)}, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if code != http.StatusOK {
		log.Fatalln(code)
	}

	err = json.Unmarshal(body, &ntp)
	if err != nil {
		log.Fatalln(err.Error())
	}

	ntp.ClientReceptionTime = time.Now().UnixNano() / int64(time.Millisecond)

	var delta = ((ntp.SrvReceptionTime - ntp.ClientTransmissionTime) +
		(ntp.SrvTransmissionTime - ntp.ClientReceptionTime)) / 2

	return delta
}

type Tag struct {
	Name string `json:"name"`
	TimestampBegin int64 `json:"timestampBegin"`
	TimestampEnd int64 `json:"timestampEnd"`
	ProductName string `json:"productName"`
	TagName string `json:"tagName"`
}

func (aftg *Connector) AddTag(tag Tag, clockDelta int64) {
	fmt.Println("Creating Tag")
	bodyBytes, err := json.Marshal(tag)
	if err != nil {
		log.Fatalln(err.Error())
	}

	code, body, err := aftg.runAftgRequest(
		"POST",
		"tags",
		bytes.NewReader(bodyBytes),
		map[string]string{"clientTransmissionTime": strconv.FormatInt(time.Now().UnixNano() / int64(time.Millisecond), 10)},
		map[string]string{"X-CLOCK-DELTA": strconv.FormatInt(clockDelta, 10)},
	)

	if err != nil {
		log.Println("Aftg Request Error", err.Error())
	}
	if code != http.StatusCreated {
		log.Println("Unexpected return code", code)
		println(string(body))
	}
}

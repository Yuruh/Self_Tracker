package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var mu sync.Mutex
var initialized uint32
var instance *AftgConnector

type AftgConnector struct {}


func GetAftgConnector() *AftgConnector {
	if atomic.LoadUint32(&initialized) == 1 {
		return instance
	}
	mu.Lock()
	defer mu.Unlock()

	if initialized == 0 {
		instance = &AftgConnector {
		}
		atomic.StoreUint32(&initialized, 1)
	}

	return instance
}

func runAftgRequest(method string, path string, requestBody, queryParams map[string]string) (int, []byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, "http://localhost:8080/" + path, nil)
	if err != nil {
		return -1, nil, err
	}
	req.Header.Add("X-API-KEY", os.Getenv("AFTG_API_KEY"))
	req.Header.Add("Content-Type", "application/json")

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

func (aftg *AftgConnector) getSrvDelay() int64 {
	var ntp NTP

	code, body, err := runAftgRequest("GET", "ntp", nil ,
		map[string]string{"clientTransmissionTime": strconv.FormatInt(time.Now().UnixNano() / int64(time.Millisecond), 10)})
	if err != nil || code != http.StatusOK {
		log.Fatalln(err.Error(), code)
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

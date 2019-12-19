package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func runTicker() {
	var ticker *time.Ticker = time.NewTicker(time.Second * 3)
	for {
		select {
		case <- ticker.C:
			println("TICK")
		}
	}
}

func buildSpotifyAuthUri() string {
	req, err := http.NewRequest("GET", "https://accounts.spotify.com/authorize", nil)

	if err != nil {
		log.Fatalln(err.Error())
	}

	query := req.URL.Query()
	query.Add("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	query.Add("response_type", "code",)
	query.Add("redirect_uri", "https://ea4f5723.ngrok.io/auth-spotify")
	query.Add("scope", "user-read-playback-state")
	req.URL.RawQuery = query.Encode()

	return req.URL.String()
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func getAccessFromRefresh() TokenResponse {
	client := &http.Client{}

	requestBody := url.Values{}
	requestBody.Set("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	requestBody.Set("client_secret", os.Getenv("SPOTIFY_CLIENT_SECRET"))
	requestBody.Set("grant_type", "refresh_token")
	requestBody.Set("refresh_token", os.Getenv("SPOTIFY_REFRESH_TOKEN"))

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(requestBody.Encode()))

	if err != nil {
		log.Fatalln(err.Error())
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(requestBody.Encode())))


	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}

	println(resp.StatusCode)
	println("response: ", string(body))

	var data TokenResponse

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return data
}

func getSpotifyTokens() {
	client := &http.Client{}

	requestBody := url.Values{}
	requestBody.Set("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	requestBody.Set("client_secret", os.Getenv("SPOTIFY_CLIENT_SECRET"))
	requestBody.Set("grant_type", "authorization_code")
	requestBody.Set("code", "<The returned code>")
	requestBody.Set("redirect_uri", "https://ea4f5723.ngrok.io/auth-spotify")

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(requestBody.Encode()))

	if err != nil {
		log.Fatalln(err.Error())
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(requestBody.Encode())))


	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}

	println(resp.StatusCode)
	println("response: ", string(body))
}

func getCurrentTrack(accessToken string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player", nil)

	if err != nil {
		log.Fatalln(err.Error())
	}
	req.Header.Add("Authorization", "Bearer " + accessToken)

	println(req.Header.Get("Authorization"))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}

	println(resp.StatusCode)
	println("response: ", string(body))
}


type NTP struct {
	SrvReceptionTime int64 `json:"srvReceptionTime"`
	ClientTransmissionTime int64 `json:"clientTransmissionTime"`
	SrvTransmissionTime int64 `json:"srvTransmissionTime"`
	ClientReceptionTime int64 `json:"clientReceptionTime"`
}

func getAftgApiSyncDelta() int64 {
	client := &http.Client{}
	var ntp NTP

	req, err := http.NewRequest("GET", "http://localhost:8080/ntp", nil)

	if err != nil {
		log.Fatalln(err.Error())
	}
	query := req.URL.Query()
	query.Add("clientTransmissionTime", strconv.FormatInt(time.Now().UnixNano() / int64(time.Millisecond), 10))
	req.URL.RawQuery = query.Encode()
	req.Header.Add("X-API-KEY", os.Getenv("AFTG_API_KEY"))

	resp, err := client.Do(req)
	ntp.ClientReceptionTime = time.Now().UnixNano() / int64(time.Millisecond)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = json.Unmarshal(body, &ntp)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var delta = ((ntp.SrvReceptionTime - ntp.ClientTransmissionTime) +
		(ntp.SrvTransmissionTime - ntp.ClientReceptionTime)) / 2

	return delta
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// runTicker()
	//	testSpotify()
	println(buildSpotifyAuthUri())
	//	getSpotifyTokens()

	// var access = getAccessFromRefresh()
	// getCurrentTrack(access.AccessToken)

	var delta = GetAftgConnector().getSrvDelay() //getAftgApiSyncDelta()

//	var roundTrip = (ntp.ClientReceptionTime - ntp.ClientTransmissionTime) - (ntp.SrvTransmissionTime - ntp.SrvReceptionTime)

	println("delta =", delta)

}
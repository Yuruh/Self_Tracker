package main

import (
	"encoding/json"
	"fmt"
	"github.com/Yuruh/Self_Tracker/aftg"

	//	"github.com/Yuruh/Self_Tracker/aftg"
	"github.com/Yuruh/Self_Tracker/spotify"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func Dummy() int64 {
	return 1
}

/*
	DESIGN IDEAS

	I want everyone to be able to use this API.

	That means we need to store users, to ask user for AFTG API Key and for Spotify OAuth

 */

func runTicker() {
//	var ticker *time.Ticker = time.NewTicker(time.Minute * 3)
	const tickInterval time.Duration = time.Second * 45
	var ticker *time.Ticker = time.NewTicker(tickInterval)

	var savedPlayer spotify.Player
	for {
		select {
		case <- ticker.C:
			var delay = aftg.GetConnector().GetSrvDelay()
			fmt.Println(delay)
			spotifyPlayer, err := spotify.GetConnector().GetCurrentTrack()
			if err != nil {
				if err, ok := err.(*spotify.TrackError); ok {
					if err.Code == spotify.NotPlaying {
						println("Not currently playing")
					} else {
						log.Fatal(err.Error())
					}
				}
			}
			if savedPlayer.Item.Id != spotifyPlayer.Item.Id {
				println("Track Changed")

				if savedPlayer.Item.Id != "" {
					var durationPlayedBeforeNextTrack time.Duration = tickInterval - time.Duration(spotifyPlayer.ProgressMs) * time.Millisecond

					var trackEndTime = time.Now().UnixNano() / int64(time.Millisecond) - spotifyPlayer.ProgressMs

					var trackBeginTime = trackEndTime -
						savedPlayer.ProgressMs -
						int64(durationPlayedBeforeNextTrack / time.Millisecond)

					fmt.Print("Title \"", savedPlayer.Item.Name, "\" played from ", time.Unix(trackBeginTime / 1000, 0))
					fmt.Println(" to", time.Unix(trackEndTime / 1000, 0))
					/*aftg.GetConnector().AddTag(aftg.Tag{
						TimestampBegin: trackBeginTime,
						TimestampEnd: trackEndTime,
						Name: savedPlayer.Item.Artists[0].Name + "_" + savedPlayer.Item.Name,
						ProductName: savedPlayer.Item.Artists[0].Name,
						TagName: savedPlayer.Item.Name,
					}, delay)*/
				}
			}
			savedPlayer = spotifyPlayer//.Copy()

//			print("Title ", spotifyPlayer.Item.Name)
//			print(" from Album ", spotifyPlayer.Item.Album.Name)
//			print(" by Artist ", spotifyPlayer.Item.Artists[0].Name)
//			println(" is playing since ", spotifyPlayer.ProgressMs, " milliseconds")
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

/*
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
*/
func main() {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}

	runTicker()

	//	println(buildSpotifyAuthUri())

	//var access = getAccessFromRefresh()
	//getCurrentTrack(access.AccessToken)



	//var delta = aftg.GetConnector().GetSrvDelay() //getAftgApiSyncDelta()

//	var roundTrip = (ntp.ClientReceptionTime - ntp.ClientTransmissionTime) - (ntp.SrvTransmissionTime - ntp.SrvReceptionTime)

	//println("delta =", delta)

}
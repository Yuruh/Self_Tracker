package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Connector struct {
	AccessToken string
	RefreshToken string
	RetryAmount int8
}

func (spotify *Connector) getAccessFromRefresh() {
	client := &http.Client{}

	requestBody := url.Values{}
	requestBody.Set("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	requestBody.Set("client_secret", os.Getenv("SPOTIFY_CLIENT_SECRET"))
	requestBody.Set("grant_type", "refresh_token")
	requestBody.Set("refresh_token", spotify.RefreshToken)

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


	var data TokenResponse

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Updating access token")
	spotify.AccessToken = data.AccessToken
}

type ErrorCode int8

const (
	NotPlaying ErrorCode = 1
)

type TrackError struct {
	err string
	Code ErrorCode
}

func (err *TrackError) Error() string {
	return "track error"
}

func (spotify *Connector) getCurrentTrack(retryAmount int8) (Player, error) {
	var spotifyPlayer Player
	if retryAmount <= 0 {
		return spotifyPlayer, errors.New("no more retry amount can't get current track")
	}
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player", nil)

	if err != nil {
		log.Fatalln(err.Error())
	}
	req.Header.Add("Authorization", "Bearer " + spotify.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Retrying...", retryAmount)
		return spotify.getCurrentTrack(retryAmount - 1)
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusOK:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err.Error())
		}

		err = json.Unmarshal(body, &spotifyPlayer)
		if err != nil {
			log.Fatalln(err.Error())
		}
		return spotifyPlayer, nil
	case resp.StatusCode == http.StatusNoContent:
//		println("no song is currently playing")
		return spotifyPlayer, &TrackError{Code: NotPlaying} // errors.New("user is not currently playing music")
	case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusBadRequest:
		println("Invalid request")
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err.Error())
		}
		println(string(body))

		spotify.getAccessFromRefresh()
		return spotify.getCurrentTrack(retryAmount - 1)
	default:
		return spotifyPlayer, errors.New("unhandled http status")
	}
}

func (spotify *Connector) GetCurrentTrack() (Player, error) {
	return spotify.getCurrentTrack(spotify.RetryAmount)
}

func BuildAuthUri(userId uint) string {
	req, err := http.NewRequest("GET", "https://accounts.spotify.com/authorize", nil)

	if err != nil {
		log.Fatalln(err.Error())
	}

	query := req.URL.Query()
	query.Add("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	query.Add("response_type", "code",)
	query.Add("redirect_uri", "http://localhost:3000/spotify-auth")
	query.Add("scope", "user-read-playback-state")
	query.Add("show_dialog", "true")
	query.Add("state", strconv.FormatUint(uint64(userId), 10))
	req.URL.RawQuery = query.Encode()

	return req.URL.String()
}

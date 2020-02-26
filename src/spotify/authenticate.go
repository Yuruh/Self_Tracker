package spotify

import (
	"encoding/json"
	"fmt"
	"github.com/Yuruh/Self_Tracker/src/database"
	"github.com/Yuruh/Self_Tracker/src/database/models"
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func ReadBody(body io.ReadCloser) string {
	defer body.Close()

	result, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return string(result)
}

func RegisterRefreshToken(context echo.Context) error {
	var user models.User = context.Get("user").(models.User)

	body := ReadBody(context.Request().Body)
	println(body)

	var data RegisterTokenRequest

	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		println(err.Error())
		return context.NoContent(http.StatusBadRequest)
	}
	state, err := strconv.Atoi(data.State)
	if err != nil || state != int(user.ID) {
		println("Bad state")
		return context.NoContent(http.StatusBadRequest)
	}
	response := getSpotifyTokens(data.Code)
	
	var api models.ApiAccess
	result := database.GetDB().Model(&user).Related(&api)
	api.Spotify = response.RefreshToken
	if result.RecordNotFound() {
		fmt.Println("Could not find matching doc, creating")
		api.UserID = user.ID
		database.GetDB().Create(&api)
	} else {
		database.GetDB().Save(&api)
	}

	return context.NoContent(http.StatusOK)
}

func getSpotifyTokens(code string) TokenResponse {
	client := &http.Client{}

	requestBody := url.Values{}
	requestBody.Set("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	requestBody.Set("client_secret", os.Getenv("SPOTIFY_CLIENT_SECRET"))
	requestBody.Set("grant_type", "authorization_code")
	requestBody.Set("code", code)
	requestBody.Set("redirect_uri", "http://localhost:3000/spotify-auth")

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

	var response TokenResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return response
}
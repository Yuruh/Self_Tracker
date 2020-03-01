package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Yuruh/Self_Tracker/src/database"
	"github.com/Yuruh/Self_Tracker/src/database/models"
	"github.com/Yuruh/Self_Tracker/src/utils"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func GetAuthUrl(c echo.Context) error {
	var user models.User = c.Get("user").(models.User)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"url": BuildAuthUri(user.ID),
		"tmp POC": c.Get("user").(models.User).Password,
	})
}

func RegisterRefreshToken(context echo.Context) error {
	var user models.User = context.Get("user").(models.User)

	body := utils.ReadBody(context.Request().Body)

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
	response, err := getSpotifyTokens(data.Code)

	if err != nil {
		return context.NoContent(http.StatusBadRequest)
	}

	var spotifyConnector models.Connector
	result := database.GetDB().Model(&user).
		Where("name = ?", "Spotify").
		Related(&spotifyConnector)
	spotifyConnector.Key = response.RefreshToken
	spotifyConnector.Registered = true
	spotifyConnector.Name = "Spotify"
	spotifyConnector.UserID = user.ID

	if result.RecordNotFound() {
		fmt.Println("Could not find matching doc, creating")
//		spotifyConnector.UserID = user.ID
		database.GetDB().Create(&spotifyConnector)
	} else {
		database.GetDB().Save(&spotifyConnector)
	}

	return context.NoContent(http.StatusOK)
}

func getSpotifyTokens(code string) (TokenResponse, error) {
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

	if resp.StatusCode != http.StatusOK {
		return response, errors.New("could not validate spotify request")
	}

	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return response, nil
}
package main

import (
	"encoding/json"
	"fmt"
	"github.com/Yuruh/Self_Tracker/src/aftg"
	"github.com/Yuruh/Self_Tracker/src/database"
	"github.com/Yuruh/Self_Tracker/src/spotify"
	_ "github.com/lib/pq"
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

func processLastPlayedSong(savedPlayer* spotify.Player, tickInterval time.Duration) {
	var delay = aftg.GetConnector().GetSrvDelay()
	fmt.Println("srv delay:", delay)
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
			aftg.GetConnector().AddTag(aftg.Tag{
				TimestampBegin: trackBeginTime,
				TimestampEnd: trackEndTime,
				Name: savedPlayer.Item.Artists[0].Name + "_" + savedPlayer.Item.Name,
				ProductName: savedPlayer.Item.Artists[0].Name,
				TagName: savedPlayer.Item.Name,
			}, delay)
		}
	}
	*savedPlayer = spotifyPlayer//.Copy()

	//			print("Title ", spotifyPlayer.Item.Name)
	//			print(" from Album ", spotifyPlayer.Item.Album.Name)
	//			print(" by Artist ", spotifyPlayer.Item.Artists[0].Name)
	//			println(" is playing since ", spotifyPlayer.ProgressMs, " milliseconds")
}

func runTicker() {
//	var ticker *time.Ticker = time.NewTicker(time.Minute * 3)
	const tickInterval time.Duration = time.Second * 45
	var ticker *time.Ticker = time.NewTicker(tickInterval)

	var savedPlayer spotify.Player
	processLastPlayedSong(&savedPlayer, tickInterval)

	for {
		select {
		case <-ticker.C:
			processLastPlayedSong(&savedPlayer, tickInterval)
		}
	}
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

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func main() {
	//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}

	// runTicker()

//	http.Handle("/foo", fooHandler)

	db := database.Connect()

	defer db.Close()

	database.RunMigration(db)

	// Create
/*	db.Create(&models.Product{Code: "L1212", Price: 1000})

	// Read
	var product Product
	db.First(&product, 1) // find product with id 1
	db.First(&product, "code = ?", "L1212") // find product with code l1212

	println(product.Code)

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	db.Delete(&product)*/

	/*_, err := sql.Open("postgres", "user=pqgotest dbname=pqgotest")
	if err != nil {
		log.Fatal(err)
	} else {
		println("Connected to db, apparently")
	}*/

	http.HandleFunc("/spotify", func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(spotify.BuildAuthUri()))
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

	})

	log.Fatal(http.ListenAndServe(":8090", nil))


	//var access = getAccessFromRefresh()
	//getCurrentTrack(access.AccessToken)



	//var delta = aftg.GetConnector().GetSrvDelay() //getAftgApiSyncDelta()

//	var roundTrip = (ntp.ClientReceptionTime - ntp.ClientTransmissionTime) - (ntp.SrvTransmissionTime - ntp.SrvReceptionTime)

	//println("delta =", delta)

}
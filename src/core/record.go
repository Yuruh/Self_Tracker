package core

import (
	"fmt"
	"github.com/Yuruh/Self_Tracker/src/aftg"
	"github.com/Yuruh/Self_Tracker/src/database/models"
	"github.com/Yuruh/Self_Tracker/src/spotify"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

func RecordActivity(context echo.Context) error {
	var user models.User = context.Get("user").(models.User)
	/*if database.GetDB().Model(&user).Related(&api).RecordNotFound() {
		log.Fatalln("Could not find Api access")
	}*/
	var spotifyConnection = spotify.Connector{RefreshToken: user.Spotify.Key, RetryAmount: 2}
	var aftgConnection = aftg.Connector{ApiKey: user.AffectTag.Key, RetryAmount: 2}
	println("Starting to record")

	go runTicker(spotifyConnection, aftgConnection)
	return context.NoContent(http.StatusOK)
}

func processLastPlayedSong(savedPlayer* spotify.Player,
	tickInterval time.Duration,
	connector *spotify.Connector,
	aftgConnector aftg.Connector,
	) {
	var delay int64 = aftgConnector.GetSrvDelay()
	fmt.Println("srv delay:", delay)

	spotifyPlayer, err := connector.GetCurrentTrack()
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
			aftgConnector.AddTag(aftg.Tag{
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

func runTicker(connector spotify.Connector, aftgConnector aftg.Connector) {
	//	var ticker *time.Ticker = time.NewTicker(time.Minute * 3)
	const tickInterval time.Duration = time.Second * 45
	var ticker *time.Ticker = time.NewTicker(tickInterval)

	var savedPlayer spotify.Player
	processLastPlayedSong(&savedPlayer, tickInterval, &connector, aftgConnector)

	for {
		select {
		case <-ticker.C:
			processLastPlayedSong(&savedPlayer, tickInterval, &connector, aftgConnector)
		}
	}
}

package core

import (
	"fmt"
	"github.com/Yuruh/Self_Tracker/src/aftg"
	"github.com/Yuruh/Self_Tracker/src/spotify"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

func RecordActivity(context echo.Context) error {
	return context.NoContent(http.StatusNotImplemented)
}

func RegisterApiKey(context echo.Context) error {
	return context.NoContent(http.StatusNotImplemented)
}

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

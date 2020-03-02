package core

import (
	"fmt"
	"github.com/Yuruh/Self_Tracker/src/aftg"
	"github.com/Yuruh/Self_Tracker/src/database"
	"github.com/Yuruh/Self_Tracker/src/database/models"
	"github.com/Yuruh/Self_Tracker/src/spotify"
	"log"
	"time"
)

type userSetup struct {
	player spotify.Player
	spotifyCon spotify.Connector
	aftgCon aftg.Connector
	ID uint
}

func findIndex(a []userSetup, x uint) int {
	for i, n := range a {
		if x == n.ID {
			return i
		}
	}
	return len(a)
}

// I Cannot use a map as i need to retrieve the address of elements

func RecordActivity() {
	const tickInterval time.Duration = time.Second * 45
	var ticker *time.Ticker = time.NewTicker(tickInterval)

	// I Cannot use a map as i need to retrieve the address of elements

	var userSetups = make([]userSetup, 0)

	println("Initializing ticker")

	// need a dynamic map of user -> savedPlayer
	// cleanup the map on every round

	for {
		select {
		case <-ticker.C:
			println("Running tick")
			var recordingUsers []models.User

			err := database.GetDB().Where("recording = ?", true).Preload("Connectors").Find(&recordingUsers)
			if err.RecordNotFound() {
				println("No recording user found")
				continue
			}

			for _, user := range recordingUsers {
				var idx = findIndex(userSetups, user.ID)
				if idx == len(userSetups) {
					println("Setting up user", user.ID)

					var spotifyConnector models.Connector
					var aftgConnector models.Connector

					for _, elem := range user.Connectors {
						println(elem.Name)
						if elem.Name == "Affect-tag" {
							aftgConnector = elem
						} else if elem.Name == "Spotify" {
							spotifyConnector = elem
						}
					}

					if !spotifyConnector.Registered || ! aftgConnector.Registered {
						continue
					}

					userSetups = append(userSetups, userSetup{
						player: spotify.Player{},
						aftgCon:aftg.Connector{ApiKey: aftgConnector.Key, RetryAmount: 2},
						spotifyCon:spotify.Connector{RefreshToken: spotifyConnector.Key, RetryAmount: 2},
						ID: user.ID,
					})
				}

//				var savedPlayer spotify.Player

				processLastPlayedSong(
					&userSetups[idx].player,
					tickInterval,
					&userSetups[idx].spotifyCon,
					userSetups[idx].aftgCon,
				)
			}

		}
	}


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
				log.Fatalln(err.Error())
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

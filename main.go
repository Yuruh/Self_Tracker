package main

import (
	"github.com/Yuruh/Self_Tracker/src/api"
	"github.com/Yuruh/Self_Tracker/src/core"
	"github.com/Yuruh/Self_Tracker/src/database"
	_ "github.com/lib/pq"
)

func Dummy() int64 {
	return 1
}


func main() {
	//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}

	// runTicker()

//	http.Handle("/foo", fooHandler)

//	db := database.Connect()

	defer database.GetDB().Close()

	database.RunMigration()

	go core.RecordActivity()

	api.RunHttpServer()


	//var access = getAccessFromRefresh()
	//getCurrentTrack(access.AccessToken)



	//var delta = aftg.GetConnector().GetSrvDelay() //getAftgApiSyncDelta()

//	var roundTrip = (ntp.ClientReceptionTime - ntp.ClientTransmissionTime) - (ntp.SrvTransmissionTime - ntp.SrvReceptionTime)

	//println("delta =", delta)

}
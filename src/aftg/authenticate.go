package aftg

import (
	"encoding/json"
	"fmt"
	"github.com/Yuruh/Self_Tracker/src/database"
	"github.com/Yuruh/Self_Tracker/src/database/models"
	"github.com/Yuruh/Self_Tracker/src/utils"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func ensureKeyValid(key string) error {
	var connector = Connector{ApiKey:key, RetryAmount:2}

	return connector.GetMe()
}

func RegisterApiKey(context echo.Context) error {
	var user models.User = context.Get("user").(models.User)

	body := utils.ReadBody(context.Request().Body)

	var data RegisterApiKeyRequest

	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		println(err.Error())
		return context.NoContent(http.StatusBadRequest)
	}

	err = ensureKeyValid(data.Key)
	if err != nil {
		println(err.Error())
		return context.String(http.StatusBadRequest, "Could not use API Key")
	}


	var aftgConnector models.Connector

	result := database.GetDB().Model(&user).
		Where("name = ?", "Affect-tag").
		Related(&aftgConnector)

	if result.Error != nil && !result.RecordNotFound() {
		log.Fatalln(result.Error.Error())
	}

	aftgConnector.Name = "Affect-tag"
	aftgConnector.Key = data.Key
	aftgConnector.Registered = true
	aftgConnector.UserID = user.ID

	if result.RecordNotFound() {
		fmt.Println("Could not find matching doc, creating")
		database.GetDB().Create(&aftgConnector)
	} else {
		fmt.Println("Matching doc already exists")
		database.GetDB().Save(&aftgConnector)
	}


	return context.JSON(http.StatusOK, aftgConnector)
}
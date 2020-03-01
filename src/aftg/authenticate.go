package aftg

import (
	"encoding/json"
	"fmt"
	"github.com/Yuruh/Self_Tracker/src/database"
	"github.com/Yuruh/Self_Tracker/src/database/models"
	"github.com/Yuruh/Self_Tracker/src/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterApiKey(context echo.Context) error {
	var user models.User = context.Get("user").(models.User)

	body := utils.ReadBody(context.Request().Body)

	var data RegisterApiKeyRequest

	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		println(err.Error())
		return context.NoContent(http.StatusBadRequest)
	}

	var aftgConnector models.Connector
	result := database.GetDB().Model(&user).Related(&aftgConnector)
	aftgConnector.Key = data.Key
	aftgConnector.Registered = true
	if result.RecordNotFound() {
		fmt.Println("Could not find matching doc, creating")
		aftgConnector.UserID = user.ID
		database.GetDB().Create(&aftgConnector)
	} else {
		database.GetDB().Save(&aftgConnector)
	}

	return context.NoContent(http.StatusOK)
}
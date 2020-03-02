package api

import (
	"encoding/json"
	"github.com/Yuruh/Self_Tracker/src/aftg"
	"github.com/Yuruh/Self_Tracker/src/database"
	"github.com/Yuruh/Self_Tracker/src/database/models"
	"github.com/Yuruh/Self_Tracker/src/spotify"
	"github.com/Yuruh/Self_Tracker/src/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func RunHttpServer()  {
	// Echo instance
	app := echo.New()
	app.HideBanner = true

	// Middleware
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Use(middleware.CORS())

	app.POST("/login", login)
	app.POST("/register", register)

	unprotectedPaths := [2]string{"/login", "/register"}

	// According to https://echo.labstack.com/middleware, "Middleware registered using Echo#Use() is only executed for paths which are registered after Echo#Use() has been called."
	// But it doesn't behave that way so for now we'll skip specific routes
	app.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims: &TokenClaims{},
		SigningKey: []byte(os.Getenv("ACCESS_TOKEN_SECRET")),
		SigningMethod: "HS256",
		ContextKey: "token",
		Skipper: func(context echo.Context) bool {
			if utils.ContainsString(unprotectedPaths[:], context.Path()) {
				return true
			}
			return false
		},
		SuccessHandler: func(context echo.Context) {
			context.Set("user", context.Get("token").(*jwt.Token).Claims.(*TokenClaims).User)
		},
	}))

	app.GET("/connectors", func (c echo.Context) error {
		var user models.User = c.Get("user").(models.User)
		var connectors []models.Connector
		database.GetDB().Where("user_id = ?", user.ID).Find(&connectors)

//		database.GetDB().Set("gorm:auto_preload", true).First(&user)

		return c.JSON(http.StatusOK, connectors)
	})

	// Routes
	app.GET("/spotify/url", spotify.GetAuthUrl)
	app.POST("/spotify/register", spotify.RegisterRefreshToken)

	app.POST("/aftg/register", aftg.RegisterApiKey)

	app.PUT("/record-activity", func(context echo.Context) error {
		var user models.User = context.Get("user").(models.User)
		body := utils.ReadBody(context.Request().Body)

		var data RecordActivityRequest

		err := json.Unmarshal([]byte(body), &data)
		if err != nil {
			println(err.Error())
			return context.NoContent(http.StatusBadRequest)
		}
		user.Recording = data.Enabled
		database.GetDB().Save(&user)
		return context.JSON(http.StatusOK, user)
	})

	// Start server
	app.Logger.Fatal(app.Start(":8090"))
}

type RecordActivityRequest struct {
	Enabled bool `json:"enabled"`
}

func (c TokenClaims) Valid() error {
	return c.StandardClaims.Valid()
}

type TokenClaims struct {
	User models.User `json:"user"`
	jwt.StandardClaims
}

func login(context echo.Context) error {
	body, err := ioutil.ReadAll(context.Request().Body)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var parsedBody models.User

	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		log.Println(err.Error())
		return context.String(http.StatusBadRequest, "")
	}
	var user models.User
	database.GetDB().Where("email = ?", parsedBody.Email).First(&user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(parsedBody.Password))

	if err != nil {
		println("User not found: " + err.Error())
		return context.String(http.StatusNotFound, "User not found")
	} else {
		user.Password = ""
		claims := &TokenClaims{
			user,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + int64(time.Hour * 24),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, _ := token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))

		return context.JSON(http.StatusOK, map[string]interface{}{"token": ss, "user": user})
	}
}

func register(context echo.Context) error {
	body, err := ioutil.ReadAll(context.Request().Body)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var parsedBody models.User

	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		log.Println(err.Error())
		return context.String(http.StatusBadRequest, "")
	}

	// todo validate input

	hash, err := bcrypt.GenerateFromPassword([]byte(parsedBody.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return context.String(http.StatusBadRequest, "Could not hash password")
	}

	clonedDb := database.GetDB().Create(&models.User{Email: parsedBody.Email, Password: string(hash)})

	if clonedDb.Error != nil {
		println(clonedDb.Error.Error())
		return context.NoContent(http.StatusBadRequest)
	}

	return context.String(http.StatusOK, "")
}

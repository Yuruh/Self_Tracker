package api

import (
	"encoding/json"
	"github.com/Yuruh/Self_Tracker/src/database"
	"github.com/Yuruh/Self_Tracker/src/database/models"
	"github.com/Yuruh/Self_Tracker/src/spotify"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
)

func messageHandler(message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(message))
	})
}

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

/*	app.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(os.Getenv("ACCESS_TOKEN_SECRET")),
	}))*/

	// Routes
	app.GET("/spotify", func(c echo.Context) error {
		return c.String(http.StatusOK, spotify.BuildAuthUri())
	})

	// Start server
	app.Logger.Fatal(app.Start(":8090"))

}

// Handler
func hello(c echo.Context) error {
	panic("lol")
	return c.String(http.StatusOK, "Hello, World!")
}


func login(context echo.Context) error {
//	var user models.User
//	db.First(&user, 1)
//	bcrypt.CompareHashAndPassword()

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
		println("ça match pas" + err.Error())
		return context.JSON(http.StatusNotFound, nil)
	} else {
		return context.JSON(http.StatusOK, nil)
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

	database.GetDB().Create(&models.User{Email: parsedBody.Email, Password: string(hash)})

	return context.String(http.StatusOK, "")
}

package api

import (
	"encoding/json"
	"github.com/Yuruh/Self_Tracker/src/database"
	"github.com/Yuruh/Self_Tracker/src/database/models"
	"github.com/Yuruh/Self_Tracker/src/spotify"
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

func ContainsString(src []string, value string) bool {
	for _, elem := range src {
		if elem == value {
			return true
		}
	}
	return false
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

	unprotectedPaths := [2]string{"/login", "/register"}

	// According to https://echo.labstack.com/middleware, "Middleware registered using Echo#Use() is only executed for paths which are registered after Echo#Use() has been called."
	// But it doesn't behave that way so for now we'll skip specific routes
	app.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims: &TokenClaims{},
		SigningKey: []byte(os.Getenv("ACCESS_TOKEN_SECRET")),
		SigningMethod: "HS256",
		Skipper: func(context echo.Context) bool {
			if ContainsString(unprotectedPaths[:], context.Path()) {
				return true
			}
			return false
		},
	}))

	// Routes
	app.GET("/spotify", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"url": spotify.BuildAuthUri(),
			"tmp POC": c.Get("user").(*jwt.Token).Claims.(*TokenClaims).User.Email},
		)
	})

	// Start server
	app.Logger.Fatal(app.Start(":8090"))

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
		println("ça match pas" + err.Error())
		return context.String(http.StatusNotFound, "User not found")
	} else {
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

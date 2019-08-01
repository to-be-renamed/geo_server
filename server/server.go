package main

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/gin-gonic/gin"
	"github.com/to-be-renamed/geo_server"
	"github.com/to-be-renamed/geo_server/generated"
	"github.com/to-be-renamed/geo_server/server/auth"
	"log"
	"time"
)

func graphqlHandler() gin.HandlerFunc {
	h := handler.GraphQL(generated.NewExecutableSchema(generated.Config{Resolvers: &geo_server.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := handler.Playground("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

type signInData struct {
	AccessCode string `json:"AccessCode" binding:"required"`
}

func signInHandler(googlePeople *auth.GooglePeople, auth *auth.AuthToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data signInData
		err := c.BindJSON(&data)
		if err != nil {
			return
		}

		err = googlePeople.Me(data.AccessCode)
		if err != nil {
			c.AbortWithStatus(400)
			return
		}

		tokenString, err := auth.TokenStringForUser("andrewmthomas87@gmail.com")
		if err != nil {
			c.AbortWithStatus(400)
			return
		}

		c.SetCookie("to-be-renamed-auth", tokenString, int(24*time.Hour/time.Second), "/", "", false, false)
	}
}

func parseJWTHandler(auth *auth.AuthToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("to-be-renamed-auth")
		if err != nil {
			c.AbortWithStatus(400)
			return
		}

		email, err := auth.UserFromTokenString(tokenString)
		if err != nil {
			c.AbortWithStatus(400)
			return
		}

		c.String(200, email)
	}
}

func main() {
	googlePeople := auth.NewGooglePeople()
	auth := auth.NewAuth("secret", "HS256")

	r := gin.Default()

	r.POST("/sign-in", signInHandler(googlePeople, auth))
	r.POST("/parse", parseJWTHandler(auth))

	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())

	log.Fatal(r.Run())
}

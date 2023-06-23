// Recipes API
//
// This is a sample recipes API.
//
//   Schemes: http
//   Host: localhost:8080
//   BasePath: /
//   Version: 1.0.0
//   Consumes:
//   - application/json
//
//   Produces:
//   - application/json
// swagger:meta
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"recipes-api/handlers"

	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipesHandler *handlers.RecipesHandler
var authHandler *handlers.AuthHandler

func init() {

	ctx := context.Background()
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	fmt.Println(status)

	recipesHandler = handlers.NewRecipeHandler(ctx, collection, redisClient)
	collectionUsers := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)

}

func main() {

	router := gin.Default()
	store, _ := redisStore.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("recipes_api", store))
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/signout", authHandler.SignOutHandler)
	router.POST("/refresh", authHandler.RefreshHandler)
	router.GET("recipes/:tag", recipesHandler.SearchRecipesHandler)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("/recipes",
			recipesHandler.NewRecipeHandler)
		authorized.PUT("/recipes/:id",
			recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id",
			recipesHandler.DeleteRecipeHandler)

	}
	router.RunTLS(":443", "certs/localhost.crt", "certs/localhost.key")

}

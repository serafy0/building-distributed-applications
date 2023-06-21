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
	authHandler = &handlers.AuthHandler{}
}

func main() {

	router := gin.Default()
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)

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
	router.Run()
}

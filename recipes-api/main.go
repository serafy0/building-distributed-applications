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
}

func main() {

	router := gin.Default()
	router.GET("/recipes", recipesHandler.ListRecipesHandler)

	authorized := router.Group("/")
	authorized.Use(AuthMiddleware())
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-KEY") != os.Getenv("X_API_KEY") {
			c.AbortWithStatus(401)
		}
		c.Next()

	}

}

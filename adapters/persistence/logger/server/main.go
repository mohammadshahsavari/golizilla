package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func main() {
	// Get the MongoDB URI from environment variables
	mongoURI := os.Getenv("MONGO_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	mongoClient = client
	defer client.Disconnect(ctx)

	// Initialize Fiber app
	app := fiber.New()

	// Use logger middleware
	app.Use(logger.New())

	// Configure CORS middleware to allow only Grafana's origin
	app.Use(cors.New(cors.Config{
		AllowOrigins: fmt.Sprintf("http://%s:%s", os.Getenv("GF_HOST"), os.Getenv("GF_PORT")),
		AllowMethods: "GET,POST",
	}))

	// Define routes
	app.Post("/query", queryHandler)

	// Start the server
	port := os.Getenv("LOGGER_PORT")
	log.Printf("Starting API on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func queryHandler(c *fiber.Ctx) error {
	var request struct {
		Collection string   `json:"collection"`
		Pipeline   []bson.M `json:"pipeline"`
	}

	// Parse the JSON body
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid payload")
	}

	collection := mongoClient.Database("logsdb").Collection(request.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Execute the MongoDB aggregation pipeline
	cursor, err := collection.Aggregate(ctx, request.Pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Return the results as JSON
	return c.JSON(results)
}

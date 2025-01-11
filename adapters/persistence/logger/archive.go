package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"golizilla/config"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ArchiveAndDelete(cfg *config.Config) error {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", cfg.MongoDbUsername, cfg.MongoDbPassword, cfg.MongoDbHost, cfg.MongoDbPort)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("logsdb").Collection("logs")

	// Define the sorting criteria (adjust based on your schema)
	sortCriteria := bson.D{{Key: "timestamp", Value: -1}} // Replace "createdAt" with your timestamp field

	// Query for the last 1000 records
	cursor, err := collection.Find(ctx, bson.D{}, options.Find().SetSort(sortCriteria).SetLimit(1000))
	if err != nil {
		return fmt.Errorf("error fetching records: %w", err)
	}
	defer cursor.Close(ctx)

	// Read the records into a slice
	var records []bson.M
	if err := cursor.All(ctx, &records); err != nil {
		return fmt.Errorf("error reading cursor: %w", err)
	}

	if len(records) == 0 {
		log.Println("No records to archive.")
		return nil
	}

	// Save records to a BSON file
	fileName := fmt.Sprintf("%s/archived_%s.bson", cfg.MongoDbArchivePath, time.Now().Format("2006-01-02T15-04-05"))
	if err := saveToBSON(fileName, records, cfg.MongoDbArchivePath); err != nil {
		return fmt.Errorf("error saving records to BSON: %w", err)
	}

	// Extract IDs of the archived records for deletion
	var ids []interface{}
	for _, record := range records {
		if id, ok := record["_id"]; ok {
			ids = append(ids, id)
		}
	}

	// Delete the archived records
	deleteFilter := bson.M{"_id": bson.M{"$in": ids}}
	_, err = collection.DeleteMany(ctx, deleteFilter)
	if err != nil {
		return fmt.Errorf("error deleting archived records: %w", err)
	}

	log.Printf("Archived and deleted %d records.", len(records))
	return nil
}

func saveToBSON(fileName string, records []bson.M, path string) error {
	// Ensure the archive directory exists
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("error creating archive directory: %w", err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating BSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(records); err != nil {
		return fmt.Errorf("error encoding BSON data: %w", err)
	}

	return nil
}

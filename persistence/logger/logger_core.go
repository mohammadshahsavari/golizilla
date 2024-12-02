package logger

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MongoDBCore is a custom Zap Core that writes logs to MongoDB.
type MongoDBCore struct {
	client     *mongo.Client
	collection *mongo.Collection
	encoder    zapcore.Encoder
	level      zapcore.LevelEnabler
}

// NewMongoDBCore creates a new MongoDBCore instance.
func NewMongoDBCore(uri, dbName, collectionName string, level zapcore.LevelEnabler) (*MongoDBCore, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)

	// Use a JSON encoder for logs
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05Z07:00"))
	})

	return &MongoDBCore{
		client:     client,
		collection: collection,
		encoder:    zapcore.NewJSONEncoder(encoderConfig),
		level:      level,
	}, nil
}

// Enabled checks if the log level is enabled.
func (m *MongoDBCore) Enabled(level zapcore.Level) bool {
	return m.level.Enabled(level)
}

// With adds structured context to the logger.
func (m *MongoDBCore) With(fields []zapcore.Field) zapcore.Core {
	return &MongoDBCore{
		client:     m.client,
		collection: m.collection,
		encoder:    m.encoder,
		level:      m.level,
	}
}

// Check determines whether the log entry should be logged.
func (m *MongoDBCore) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if m.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, m)
	}
	return checkedEntry
}

// Write writes the log entry to MongoDB.
func (m *MongoDBCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// Encode the log entry
	buffer, err := m.encoder.EncodeEntry(entry, fields)
	if err != nil {
		return err
	}

	// Convert the log entry to a BSON document
	var logDoc bson.M
	if err := bson.UnmarshalExtJSON(buffer.Bytes(), false, &logDoc); err != nil {
		return err
	}

	// Add the log to MongoDB
	_, err = m.collection.InsertOne(context.TODO(), logDoc)
	return err
}

// Sync ensures all buffered logs are written.
func (m *MongoDBCore) Sync() error {
	return m.client.Disconnect(context.TODO())
}

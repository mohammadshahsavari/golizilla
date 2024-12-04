package logger

import (
	"context"
	"log"
	"sync"
	"time"

	"golizilla/handler/middleware" // Import middleware for session handling

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MongoDBCore is a custom Zap Core that writes logs to MongoDB.
type MongoDBCore struct {
	client       *mongo.Client
	collection   *mongo.Collection
	encoder      zapcore.Encoder
	level        zap.AtomicLevel
	fields       []zapcore.Field
	levelEnabler zapcore.LevelEnabler
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
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "level"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	return &MongoDBCore{
		client:       client,
		collection:   collection,
		encoder:      encoder,
		level:        zap.NewAtomicLevelAt(zapcore.DebugLevel),
		levelEnabler: level,
		fields:       []zapcore.Field{},
	}, nil
}

// Enabled checks if the log level is enabled.
func (m *MongoDBCore) Enabled(level zapcore.Level) bool {
	return m.levelEnabler.Enabled(level)
}

// With adds structured context to the logger.
func (m *MongoDBCore) With(fields []zapcore.Field) zapcore.Core {
	return &MongoDBCore{
		client:       m.client,
		collection:   m.collection,
		encoder:      m.encoder.Clone(),
		level:        m.level,
		levelEnabler: m.levelEnabler,
		fields:       append(m.fields, fields...),
	}
}

// Check determines whether the log entry should be logged.
func (m *MongoDBCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if m.Enabled(entry.Level) {
		return ce.AddCore(entry, m)
	}
	return ce
}

// Write writes the log entry to MongoDB.
func (m *MongoDBCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	allFields := append(m.fields, fields...)

	buffer, err := m.encoder.EncodeEntry(entry, allFields)
	if err != nil {
		return err
	}

	var logDoc bson.M
	if err := bson.UnmarshalExtJSON(buffer.Bytes(), false, &logDoc); err != nil {
		return err
	}

	_, err = m.collection.InsertOne(context.TODO(), logDoc)
	return err
}

// Sync ensures all buffered logs are written.
func (m *MongoDBCore) Sync() error {
	return m.client.Disconnect(context.TODO())
}

// Logger wraps zap.Logger and provides structured logging methods.
type Logger struct {
	zapLogger *zap.Logger
}

// LogFields defines the required fields for logging.
type LogFields struct {
	Service       string
	Endpoint      string
	UserID        string
	SessionID     string
	TransactionID string
	TraceID       string
	Message       string
	Context       interface{}
}

// addDefaultFields adds default values for trace_id, transaction_id, and session_id.
func (l *Logger) addDefaultFields(fields LogFields, c *fiber.Ctx) LogFields {
	if c != nil {
		if fields.TraceID == "" {
			fields.TraceID, _ = c.Locals("trace_id").(string)
		}
		if fields.TransactionID == "" {
			fields.TransactionID, _ = c.Locals("transaction_id").(string)
		}
		if fields.SessionID == "" {
			session, _ := middleware.Store.Get(c)
			fields.SessionID = session.ID()
		}
		if fields.UserID == "" {
			session, _ := middleware.Store.Get(c)
			fields.UserID, _ = session.Get("user_id").(string)
		}
		if fields.Endpoint == "" {
			fields.Endpoint = c.OriginalURL()
		}
	}
	return fields
}

// LogInfo logs an info-level message with default fields included.
func (l *Logger) LogInfo(fields LogFields, c *fiber.Ctx) {
	fields = l.addDefaultFields(fields, c)
	l.zapLogger.Info(fields.Message,
		zap.String("service", fields.Service),
		zap.String("endpoint", fields.Endpoint),
		zap.String("user_id", fields.UserID),
		zap.String("session_id", fields.SessionID),
		zap.String("transaction_id", fields.TransactionID),
		zap.String("trace_id", fields.TraceID),
		zap.Any("context", fields.Context),
	)
}

// LogError logs an error-level message with default fields included.
func (l *Logger) LogError(fields LogFields, c *fiber.Ctx) {
	fields = l.addDefaultFields(fields, c)
	l.zapLogger.Error(fields.Message,
		zap.String("service", fields.Service),
		zap.String("endpoint", fields.Endpoint),
		zap.String("user_id", fields.UserID),
		zap.String("session_id", fields.SessionID),
		zap.String("transaction_id", fields.TransactionID),
		zap.String("trace_id", fields.TraceID),
		zap.Any("context", fields.Context),
	)
}

// LogDebug logs a debug-level message with default fields included.
func (l *Logger) LogDebug(fields LogFields, c *fiber.Ctx) {
	fields = l.addDefaultFields(fields, c)
	l.zapLogger.Debug(fields.Message,
		zap.String("service", fields.Service),
		zap.String("endpoint", fields.Endpoint),
		zap.String("user_id", fields.UserID),
		zap.String("session_id", fields.SessionID),
		zap.String("transaction_id", fields.TransactionID),
		zap.String("trace_id", fields.TraceID),
		zap.Any("context", fields.Context),
	)
}

// LogWarning logs a warning-level message with default fields included.
func (l *Logger) LogWarning(fields LogFields, c *fiber.Ctx) {
	fields = l.addDefaultFields(fields, c)
	l.zapLogger.Warn(fields.Message,
		zap.String("service", fields.Service),
		zap.String("endpoint", fields.Endpoint),
		zap.String("user_id", fields.UserID),
		zap.String("session_id", fields.SessionID),
		zap.String("transaction_id", fields.TransactionID),
		zap.String("trace_id", fields.TraceID),
		zap.Any("context", fields.Context),
	)
}

// LogFatal logs a fatal-level message with default fields included.
func (l *Logger) LogFatal(fields LogFields, c *fiber.Ctx) {
	fields = l.addDefaultFields(fields, c)
	l.zapLogger.Fatal(fields.Message,
		zap.String("service", fields.Service),
		zap.String("endpoint", fields.Endpoint),
		zap.String("user_id", fields.UserID),
		zap.String("session_id", fields.SessionID),
		zap.String("transaction_id", fields.TransactionID),
		zap.String("trace_id", fields.TraceID),
		zap.Any("context", fields.Context),
	)
}

var (
	singletonLogger *Logger
	once            sync.Once
)

// Initialize sets up the singleton logger instance.
func Initialize(uri, dbName, collectionName string, level zapcore.LevelEnabler) error {
	var err error
	once.Do(func() {
		mongoCore, coreErr := NewMongoDBCore(uri, dbName, collectionName, level)
		if coreErr != nil {
			err = coreErr
			return
		}
		singletonLogger = &Logger{zapLogger: zap.New(mongoCore)}
	})
	return err
}

// GetLogger returns the singleton logger instance.
func GetLogger() *Logger {
	if singletonLogger == nil {
		log.Fatal("Logger is not initialized. Call Initialize() first.")
	}
	return singletonLogger
}

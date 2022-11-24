package logger

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/AlexandrBurak/TaxiApp/internal/config"
)

type Logger struct {
	collection *mongo.Collection
}

func NewLogger() (Logger, error) {
	cfg, err := config.GetAppCfg()
	if err != nil {
		return Logger{}, err
	}
	sess := cfg.MONGO_URL
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(sess))
	if err != nil {
		return Logger{}, err
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return Logger{collection: nil}, err
	}
	collection := client.Database("logDB").Collection("log")
	return Logger{collection: collection}, nil
}
func (l *Logger) Log(err error) {
	s := err.Error()

	doc := bson.D{primitive.E{Key: "level", Value: "info"}, primitive.E{Key: "time", Value: time.Now()}, primitive.E{Key: "message", Value: s}}
	_, err = l.collection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Println(err.Error())
	}
}

func (l *Logger) Error(err error) {
	s := err.Error()
	doc := bson.D{primitive.E{Key: "level", Value: "error"}, primitive.E{Key: "time", Value: time.Now()}, primitive.E{Key: "message", Value: s}}
	_, err = l.collection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Println(err.Error())
	}
}

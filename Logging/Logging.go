package Logging

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogType struct {
	Query        []byte
	Response     []byte
	ComputeTime  int64
	RegisterTime int64
}

var clientOption *options.ClientOptions

var _dbName, _dbColName string

func InitLogger(dbPath string, dbUser string, dbPass string, dbName string, dbColName string) {
	clientOption = options.Client().ApplyURI(dbPath) //.SetAuth(options.Credential{Username: dbUser, Password: dbPass}) Thros Some Error

	_dbName = dbName
	_dbColName = dbColName
}

func Log(query []byte, resp []byte, comTime int64) error {
	log := LogType{
		Query:        query,
		Response:     resp,
		ComputeTime:  comTime,
		RegisterTime: int64(time.Now().Nanosecond()),
	}

	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		return err
	}

	Logs := client.Database(_dbName).Collection(_dbColName)

	_, err = Logs.InsertOne(context.TODO(), log)
	if err != nil {
		return err
	}

	client.Disconnect(context.TODO())
	return nil
}

package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var lockInitMongoDBPool sync.Mutex

var mongoPool *mongo.Client
var connsCheckedOut int

const (
	mongoConnectTimeout         = 8 * time.Second
	mongoServerSelectionTimeout = 30 * time.Second
	mongoWriteTimeout           = 8 * time.Second
)

func (dbPool DBPool) handleMongoPoolMonitor(evt *event.PoolEvent) {
	switch evt.Type {
	case event.GetSucceeded:
		connsCheckedOut++
	case event.ConnectionReturned:
		connsCheckedOut--
	case event.PoolClosedEvent:
		connsCheckedOut = 0
	}

	if evt.PoolOptions != nil {
		dbPool.Logger.Trace("", fmt.Sprintf("Mongo DB Event: %s, Address: %s, PoolOptions: %+v, ConnsCheckedOut: %d, Reason: %s",
			evt.Type, evt.Address, *evt.PoolOptions, connsCheckedOut, evt.Reason))
	} else {
		dbPool.Logger.Trace("", fmt.Sprintf("Mongo DB Event: %s, Address: %s, PoolOptions: null, ConnsCheckedOut: %d, Reason: %s",
			evt.Type, evt.Address, connsCheckedOut, evt.Reason))
	}
}

func (dbPool DBPool) pingMongoDB(transID string) error {
	ctxPing, cancelPing := context.WithTimeout(context.Background(), mongoConnectTimeout)
	defer cancelPing()

	err := mongoPool.Ping(ctxPing, readpref.Primary())

	if err != nil {
		dbPool.Logger.Error(transID, "Can not verify a connection to Mongo DB because "+err.Error(), err)
	}

	return err
}

func (dbPool DBPool) initMongoDBPool(transID string) (*mongo.Client, error) {
	lockInitMongoDBPool.Lock()
	defer lockInitMongoDBPool.Unlock()

	var err error

	if mongoPool != nil {
		err = dbPool.pingMongoDB(transID)

		if err != nil {
			dbPool.Disconnect(transID)
		} else {
			return mongoPool, err
		}
	}

	// ApplyURI = "mongodb://userName:Password@host1:portNo,host2:portNo,host3:portNo/?replicaSet=namedOfReplicaSet"
	clientOptions := options.Client()
	clientOptions.ApplyURI(dbPool.DataSourceName)

	writeConcern := writeconcern.Majority()
	writeConcern.WTimeout = mongoWriteTimeout
	clientOptions.SetWriteConcern(writeConcern)

	clientOptions.SetReadConcern(readconcern.Majority())

	clientOptions.SetMaxPoolSize(uint64(dbPool.MaxOpenConns))
	clientOptions.SetMinPoolSize(uint64(dbPool.MaxIdleConns))
	clientOptions.SetMaxConnIdleTime(dbPool.MaxLifetime)
	clientOptions.SetServerSelectionTimeout(mongoServerSelectionTimeout)

	if len(dbPool.AppName) > 0 {
		// The server prints the appname to the MongoDB logs upon establishing the connection.
		// It is also recorded in the slow query logs and profile collections.
		clientOptions.SetAppName(dbPool.AppName)
	}

	poolMonitor := &event.PoolMonitor{
		Event: dbPool.handleMongoPoolMonitor,
	}

	clientOptions.SetPoolMonitor(poolMonitor)

	ctxConn, cancelConn := context.WithTimeout(context.Background(), mongoConnectTimeout)
	defer cancelConn()

	mongoPool, err = mongo.Connect(ctxConn, clientOptions)

	if err != nil {
		dbPool.Logger.Error(transID, "Can not connect to Mongo DB because "+err.Error())
		return mongoPool, err
	}

	err = dbPool.pingMongoDB(transID)

	return mongoPool, err
}

func (dbPool DBPool) GetMongoDBPool(transID string) (*mongo.Client, error) {
	var err error

	if mongoPool == nil {
		return dbPool.initMongoDBPool(transID)
	}

	err = dbPool.pingMongoDB(transID)

	if err != nil {
		return dbPool.initMongoDBPool(transID)
	}

	return mongoPool, err
}

func (dbPool DBPool) Disconnect(transID string) error {
	var err error

	if mongoPool != nil {
		err = mongoPool.Disconnect(context.Background())

		if err == nil {
			mongoPool = nil
			dbPool.Logger.Info(transID, "Disconnect to Mongo DB success")
		} else {
			dbPool.Logger.Error(transID, "Can not disconnect to Mongo DB because "+err.Error())
		}
	}

	return err
}

func GetCollection(dbName string, collectionName string) *mongo.Collection {
	return mongoPool.Database(dbName).Collection(collectionName)
}

// NewFindOptions pageNo start with 0
func NewFindOptions(limitRecords int64, pageNo int64) *options.FindOptions {
	findOptions := options.Find()

	if limitRecords > 0 {
		findOptions.SetLimit(limitRecords)
		findOptions.SetSkip(limitRecords * pageNo)
	}

	return findOptions
}

func NewFindOneOptions(sort interface{}) *options.FindOneOptions {
	findOneOptions := options.FindOne()

	if sort != nil {
		findOneOptions.SetSort(sort)
	}

	return findOneOptions
}

func NewTransactionOptions() *options.TransactionOptions {
	txOpts := options.Transaction()
	writeConcern := writeconcern.Majority()
	writeConcern.WTimeout = mongoWriteTimeout
	txOpts.SetWriteConcern(writeConcern)
	txOpts.SetReadConcern(readconcern.Snapshot())
	return txOpts
}

func RunCommand(dbName string, cmd interface{}, opts ...*options.RunCmdOptions) *mongo.SingleResult {
	db := mongoPool.Database(dbName)
	return db.RunCommand(context.Background(), cmd, opts...)
}

func CreateCollection(dbName string, collectionName string) *mongo.SingleResult {
	return RunCommand(dbName, bson.D{{"create", collectionName}})
}

func StringToObjectID(s string) primitive.ObjectID {
	objId, _ := primitive.ObjectIDFromHex(s)
	return objId
}

func StringToObjectIDPointer(s string) *primitive.ObjectID {
	objId, _ := primitive.ObjectIDFromHex(s)
	return &objId
}

func NewDateTimeMongo(dt time.Time) *primitive.DateTime {
	if !dt.IsZero() {
		dtMongo := primitive.NewDateTimeFromTime(dt)
		return &dtMongo
	} else {
		return nil
	}
}

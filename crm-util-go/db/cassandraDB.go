package db

import (
	"crm-util-go/validate"
	"github.com/gocql/gocql"
	"sync"
)

var lockSessionOne sync.Mutex
var lockSessionTrx sync.Mutex
var sessionOne *gocql.Session
var sessionTrx *gocql.Session

func (csdPool CassandraPool) NewCluster(consistencyLevel gocql.Consistency) *gocql.ClusterConfig {
	var emptyDbAuthen DbAuthen

	// connect to the cluster
	cluster := gocql.NewCluster(csdPool.Hosts...)
	cluster.Port = csdPool.Port
	cluster.Keyspace = csdPool.Keyspace
	cluster.Consistency = consistencyLevel
	cluster.Timeout = csdPool.ConnectTimeout
	cluster.ConnectTimeout = csdPool.ConnectTimeout

	if validate.HasStringValue(csdPool.CQLVersion) {
		cluster.CQLVersion = csdPool.CQLVersion
	}

	if csdPool.NativeProtoVersion > 0 {
		cluster.ProtoVersion = csdPool.NativeProtoVersion
	}

	if csdPool.DbAuthen != emptyDbAuthen {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: csdPool.DbAuthen.Username,
			Password: csdPool.DbAuthen.Password,
		}
	}

	return cluster
}

func (csdPool CassandraPool) initSessionOne(transID string) (*gocql.Session, error) {

	lockSessionOne.Lock()
	defer lockSessionOne.Unlock()

	var err error

	if sessionOne != nil {
		if !sessionOne.Closed() {
			return sessionOne, err
		}
	}

	cluster := csdPool.NewCluster(gocql.One)

	sessionOne, err = cluster.CreateSession()

	if err != nil {
		csdPool.Logger.Error(transID, "Can not connect to Cassandra DB because " + err.Error(), err)
		return sessionOne, err
	} else {
		csdPool.Logger.Info(transID, "Connect to Cassandra DB success")
	}

	return sessionOne, err
}

func (csdPool CassandraPool) initSessionTransaction(transID string) (*gocql.Session, error) {

	lockSessionTrx.Lock()
	defer lockSessionTrx.Unlock()

	var err error

	if sessionTrx != nil {
		if !sessionTrx.Closed() {
			return sessionTrx, err
		}
	}

	cluster := csdPool.NewCluster(gocql.LocalQuorum)

	sessionTrx, err = cluster.CreateSession()

	if err != nil {
		csdPool.Logger.Error(transID, "Can not connect to Cassandra DB because " + err.Error(), err)
		return sessionTrx, err
	} else {
		csdPool.Logger.Info(transID, "Connect to Cassandra DB success")
	}

	return sessionTrx, err
}

func (csdPool CassandraPool) GetSessionOne(transID string) (*gocql.Session, error) {
	if sessionOne == nil {
		return csdPool.initSessionOne(transID)
	} else if sessionOne.Closed() {
		return csdPool.initSessionOne(transID)
	} else {
		return sessionOne, nil
	}
}

func (csdPool CassandraPool) GetSessionTransaction(transID string) (*gocql.Session, error) {
	if sessionTrx == nil {
		return csdPool.initSessionTransaction(transID)
	} else if sessionTrx.Closed() {
		return csdPool.initSessionTransaction(transID)
	} else {
		return sessionTrx, nil
	}
}

func (csdPool CassandraPool) CloseAllSession() {
	if sessionOne != nil {
		sessionOne.Close()
	}

	if sessionTrx != nil {
		sessionTrx.Close()
	}
}
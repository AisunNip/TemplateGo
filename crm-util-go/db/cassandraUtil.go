package db

import (
	"crm-util-go/logging"
	"crm-util-go/validate"
	"errors"
	config "github.com/spf13/viper"
	"sort"
	"strings"
	"time"
)

const cassandraConnectTimeout = time.Duration(3) * time.Second

func NewCassandraPool(hosts []string, username string, password string,
	keyspace string, logger *logging.PatternLogger) CassandraPool {

	// CQLVersion: SELECT cql_version FROM system.local;
	// NativeProtoVersion: show version
	var csdPool CassandraPool
	csdPool.Hosts = hosts
	csdPool.DbAuthen.Username = username
	csdPool.DbAuthen.Password = password
	csdPool.Keyspace = keyspace

	csdPool.Port = 9042
	csdPool.ConnectTimeout = cassandraConnectTimeout
	csdPool.CQLVersion = "3.4.4"
	csdPool.NativeProtoVersion = 4
	csdPool.Logger = logger

	return csdPool
}

func (csdPool CassandraPool) GetConfigVal(transID string, key string) (string, error) {
	var val string

	session, err := csdPool.GetSessionOne(transID)

	if err != nil {
		csdPool.Logger.Error(transID, "Cassandra DB Error: "+err.Error(), err)
		return val, err
	}

	iterator := session.Query("SELECT value FROM config WHERE key = ?", key).Iter()
	scanner := iterator.Scanner()

	isFound := false

	for scanner.Next() {
		err = scanner.Scan(&val)

		if err != nil {
			return val, err
		}

		isFound = true
	}

	err = iterator.Close()

	if err != nil {
		return val, errors.New("Cassandra database error " + err.Error())
	}

	if !isFound {
		err = errors.New("Data config.key=" + key + " not found in Cassandra database")
	}

	return val, err
}

func (csdPool CassandraPool) LoadConfig(transID string, keyList []interface{}) error {
	var cqlBuilder strings.Builder
	cqlBuilder.WriteString("SELECT key, value FROM config WHERE key in (")

	for i := 0; i < len(keyList); i++ {
		if i == 0 {
			cqlBuilder.WriteString("?")
		} else {
			cqlBuilder.WriteString(",?")
		}
	}

	cqlBuilder.WriteString(")")

	session, err := csdPool.GetSessionOne(transID)

	if err != nil {
		csdPool.Logger.Error(transID, "Cassandra DB Error: "+err.Error(), err)
		return err
	}

	iter := session.Query(cqlBuilder.String(), keyList...).Iter()

	row := make(map[string]interface{})
	for iter.MapScan(row) {
		config.Set(row["key"].(string), row["value"].(string))
		row = map[string]interface{}{}
	}

	err = iter.Close()

	return err
}

type MappingBean struct {
	MappingType string `json:"mappingType"`
	FromValue   string `json:"fromValue"`
	ToValue     string `json:"toValue"`
}

func (csdPool CassandraPool) GetMapping(transID string, mappingTypeList []string, fromValue string) ([]MappingBean, error) {
	var mappingBeanList []MappingBean

	var bindValue []interface{}
	var cqlBuilder strings.Builder
	cqlBuilder.WriteString("SELECT type, from_value, to_value FROM mapping WHERE type ")

	if len(mappingTypeList) > 1 {
		cqlBuilder.WriteString("IN (")

		for i := 0; i < len(mappingTypeList); i++ {
			bindValue = append(bindValue, mappingTypeList[i])

			if i == 0 {
				cqlBuilder.WriteString("?")
			} else {
				cqlBuilder.WriteString(",?")
			}
		}

		cqlBuilder.WriteString(")")
	} else {
		bindValue = append(bindValue, mappingTypeList[0])
		cqlBuilder.WriteString("= ?")
	}

	if validate.HasStringValue(fromValue) {
		bindValue = append(bindValue, fromValue)
		cqlBuilder.WriteString(" AND from_value = ?")
	}

	if len(mappingTypeList) > 1 {
		cqlBuilder.WriteString(" ALLOW FILTERING")
	}

	session, err := csdPool.GetSessionOne(transID)

	if err != nil {
		csdPool.Logger.Error(transID, "Cassandra DB Error: "+err.Error(), err)
		return mappingBeanList, err
	}

	iterator := session.Query(cqlBuilder.String(), bindValue...).Iter()

	m := map[string]interface{}{}
	for iterator.MapScan(m) {
		mappingBean := MappingBean{
			MappingType: m["type"].(string),
			FromValue:   m["from_value"].(string),
			ToValue:     m["to_value"].(string),
		}

		mappingBeanList = append(mappingBeanList, mappingBean)

		m = map[string]interface{}{}
	}

	err = iterator.Close()

	if err != nil {
		return mappingBeanList, errors.New("Cassandra database error " + err.Error())
	}

	if len(mappingBeanList) > 0 {
		if len(mappingTypeList) > 1 {
			sort.SliceStable(mappingBeanList, func(i, j int) bool {
				return mappingBeanList[i].MappingType < mappingBeanList[j].MappingType
			})
		}
	} else {
		err = errors.New("Data mappingTypeList=" + strings.Join(mappingTypeList, ",") + " not found in Cassandra database")
	}

	return mappingBeanList, err
}

func (csdPool CassandraPool) GetMappingByType(transID string, mappingType string) ([]MappingBean, error) {
	mappingTypeList := []string{mappingType}
	return csdPool.GetMapping(transID, mappingTypeList, "")
}

func (csdPool CassandraPool) GetMappingByTypeList(transID string, mappingTypeList []string) ([]MappingBean, error) {
	return csdPool.GetMapping(transID, mappingTypeList, "")
}

func (csdPool CassandraPool) FilterMappingByType(mappingList []MappingBean, mappingType string) []MappingBean {
	var resultMappingList []MappingBean

	if mappingList != nil {
		for _, mappingBean := range mappingList {
			if mappingType == mappingBean.MappingType {
				resultMappingList = append(resultMappingList, mappingBean)
			}
		}
	}

	return resultMappingList
}

func (csdPool CassandraPool) GetMappingValue(mappingList []MappingBean, mappingType string, fromValue string) (string, error) {
	toValue := ""
	isFound := false

	if mappingList != nil {
		for _, mappingBean := range mappingList {
			if mappingType == mappingBean.MappingType && fromValue == mappingBean.FromValue {
				toValue = mappingBean.ToValue
				isFound = true
				break
			}
		}
	}

	if !isFound {
		return toValue, errors.New("Data Mapping type=" + mappingType + ", fromValue=" + fromValue + " not found in Cassandra database.")
	}

	return toValue, nil
}

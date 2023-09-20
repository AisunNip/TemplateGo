package db

import (
	"crm-util-go/common"
	"crm-util-go/logging"
	"encoding/json"
	"fmt"
	config "github.com/spf13/viper"
	"testing"
)

var csdLogger *logging.PatternLogger
var csdPool CassandraPool

func initCassandra() {
	csdLogger = logging.InitInboundLogger("crm-util-go", logging.CrmDatabase)
	csdLogger.Level = logging.LEVEL_ALL

	// ############ Init Config ############
	config.SetConfigName("uat2")
	config.SetConfigType("yaml")
	config.AddConfigPath("../config")
	err := config.ReadInConfig()

	if err != nil {
		panic("Error load a configuration file")
	}

	hosts := config.GetStringSlice("crm.cassandra.cluster")
	username := config.GetString("crm.cassandra.username")
	password := config.GetString("crm.cassandra.password")
	keyspace := config.GetString("crm.cassandra.keyspace")

	csdPool = NewCassandraPool(hosts, username, password, keyspace, csdLogger)
}

func TestCassandraGetConfigVal(t *testing.T) {
	transID := common.NewUUID()

	initCassandra()
	defer csdPool.CloseAllSession()

	// ############ GetConfigVal ############
	// "siebel.url", "siebel.header.user", "siebel.header.pass"
	val, err := csdPool.GetConfigVal(transID, "siebel.header.pass")
	csdLogger.Info(transID, "ConfigVal: "+val)

	if err != nil {
		t.Errorf("TestQueryCassandra Error %v", err.Error())
	}
}

func TestCassandraLoadConfig(t *testing.T) {
	transID := common.NewUUID()

	initCassandra()
	defer csdPool.CloseAllSession()

	keyList := []interface{}{"siebel.url", "siebel.header.user", "siebel.header.pass"}
	err := csdPool.LoadConfig(transID, keyList)

	if err != nil {
		t.Errorf("TestCassandraLoadConfig Error %v", err.Error())
	}

	csdLogger.Info(transID, "siebel.url: "+config.GetString("siebel.url"))
	csdLogger.Info(transID, "siebel.header.user: "+config.GetString("siebel.header.user"))
	csdLogger.Info(transID, "siebel.header.pass: "+config.GetString("siebel.header.pass"))
}

func TestCassandraGetMapping(t *testing.T) {
	transID := common.NewUUID()

	initCassandra()
	defer csdPool.CloseAllSession()

	// ############ GetMapping ############
	mappingBeanList, err := csdPool.GetMappingByType(transID, "TITLE")
	// mappingBeanList, err := csdPool.GetMapping(transID,"TITLE", "พล.อ.อ.")

	if err != nil {
		t.Errorf("TestQueryListCassandra Error %s", err.Error())
	} else {
		csdLogger.Info(transID, "GetMapping Success")

		for _, mappingBean := range mappingBeanList {
			csdLogger.Info(transID, fmt.Sprintf("Type: %s, FromVal: %s, ToVal: %s",
				mappingBean.MappingType, mappingBean.FromValue, mappingBean.ToValue))
		}
	}
}

func TestCassandraGetMappingList(t *testing.T) {
	transID := common.NewUUID()

	initCassandra()
	defer csdPool.CloseAllSession()

	mappingByTypeList := []string{"TITLE", "IDTYPE"}
	mappingBeanList, err := csdPool.GetMappingByTypeList(transID, mappingByTypeList)

	if err != nil {
		t.Errorf("TestQueryListCassandra Error %s", err.Error())
	} else {
		csdLogger.Info(transID, "GetMapping Success")
		binary, _ := json.MarshalIndent(mappingBeanList, "", "  ")
		fmt.Println(string(binary))
	}
}

func TestCassandraGetMappingValue(t *testing.T) {
	transID := common.NewUUID()

	initCassandra()
	defer csdPool.CloseAllSession()

	mappingByTypeList := []string{"TITLE", "IDTYPE"}
	mappingBeanList, err := csdPool.GetMappingByTypeList(transID, mappingByTypeList)

	if err != nil {
		t.Errorf("TestQueryListCassandra Error %s", err.Error())
	} else {
		csdLogger.Info(transID, "GetMapping Success")

		toValue, err := csdPool.GetMappingValue(mappingBeanList, "IDTYPE", "Personal Identity")

		if err != nil {
			csdLogger.Info(transID, err.Error())
		} else {
			csdLogger.Info(transID, "toValue="+toValue)
		}
	}
}

/*
	 err := session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`,
			"me", gocql.TimeUUID(), "hello world").WithContext(ctx).Exec()
*/

// http://www.code2succeed.com/go-cassandra-crud-example/

/*
	func getEmps() []Emp {
		fmt.Println("Getting all Employees")
		var emps []Emp
		m := map[string]interface{}{}

		iter := Session.Query("SELECT * FROM emps").Iter()
		for iter.MapScan(m) {
			emps = append(emps, Emp{
				id:        m["empid"].(string),
				firstName: m["first_name"].(string),
				lastName:  m["last_name"].(string),
				age:       m["age"].(int),
			})
			m = map[string]interface{}{}
		}

		return emps
	}
*/

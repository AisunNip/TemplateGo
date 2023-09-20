package db

import (
	"context"
	"crm-util-go/common"
	"crm-util-go/errorcode"
	"crm-util-go/logging"
	"crm-util-go/pointer"
	"crm-util-go/timeUtil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"testing"
	"time"
)

/*
	MongoDB are stored in a binary representation called BSON (Binary-encoded JSON).
	https://bsonspec.org/
	BSON Date = golang int64  Unix epoch (Jan 1, 1970) UTC datetime

	You can specify a particular date by passing an ISO-8601 date string with a year within
	the inclusive range 0 through 9999 to the new Date() constructor or the ISODate() function. These functions accept the following formats:

	new Date("<YYYY-mm-dd>") returns the ISODate with the specified date.
	new Date("<YYYY-mm-ddTHH:MM:ss>") specifies the datetime in the client's local timezone and returns the ISODate with the specified datetime in UTC.
	new Date("<YYYY-mm-ddTHH:MM:ssZ>") specifies the datetime in UTC and returns the ISODate with the specified datetime in UTC.
	new Date(<integer>) specifies the datetime as milliseconds since the UNIX epoch (Jan 1, 1970), and returns the resulting ISODate instance.

	MongoDB's data model allows related data to be stored together in a single document.
	We estimate that 80%-90% of applications that model their data in a way that leverages the document model
	will not require multi-document transactions. However, MongoDB supports multi-document ACID transactions
	for the use cases that require them. Developers appreciate the flexibility of being able to model
	their data in a way that does not typically require multi-document transactions but having
	multi-document transaction capabilities available in the event they do.
*/

const (
	mongoQueryTimeout  = time.Duration(3) * time.Second
	mongoInsertTimeout = time.Duration(3) * time.Second
	mongoUpdateTimeout = time.Duration(3) * time.Second
	mongoDeleteTimeout = time.Duration(3) * time.Second
	databaseName       = "test"
)

type EmployeeResp struct {
	Code         string        `json:"code,omitempty"`
	Msg          string        `json:"msg,omitempty"`
	TransID      string        `json:"transID,omitempty"`
	TotalRecords int64         `json:"totalRecords"`
	EmployeeList []EmployeeDep `json:"employeeList,omitempty"`
}

type Employee struct {
	RowID        *primitive.ObjectID `json:"rowId,omitempty" bson:"_id,omitempty"`
	IDNo         *string             `json:"idNo,omitempty" bson:"idNo"`
	FirstName    *string             `json:"firstName,omitempty" bson:"firstName"`
	LastName     *string             `json:"lastName,omitempty" bson:"lastName"`
	Age          *int                `json:"age,omitempty" bson:"age"`
	Birthday     *primitive.DateTime `json:"birthday,omitempty" bson:"birthday"`
	DepartmentID *primitive.ObjectID `json:"departmentID,omitempty" bson:"departmentID"`
	Address      *Address            `json:"address,omitempty" bson:"address"`
}

type EmployeeDep struct {
	RowID        *primitive.ObjectID `json:"rowId,omitempty" bson:"_id,omitempty"`
	IDNo         *string             `json:"idNo,omitempty" bson:"idNo"`
	FirstName    *string             `json:"firstName,omitempty" bson:"firstName"`
	LastName     *string             `json:"lastName,omitempty" bson:"lastName"`
	Age          *int                `json:"age,omitempty" bson:"age"`
	Birthday     *primitive.DateTime `json:"birthday,omitempty" bson:"birthday"`
	DepartmentID *primitive.ObjectID `json:"departmentID,omitempty" bson:"departmentID"`
	Address      *Address            `json:"address,omitempty" bson:"address"`
	Department   *[]Department       `json:"department,omitempty" bson:"department"`
}

type Department struct {
	RowID          *primitive.ObjectID `json:"rowId,omitempty" bson:"_id,omitempty"`
	DepartmentID   *string             `json:"departmentID,omitempty" bson:"departmentID"`
	DepartmentName *string             `json:"departmentName,omitempty" bson:"departmentName"`
	CreatedDate    *primitive.DateTime `json:"createdDate,omitempty" bson:"createdDate"`
}

type Address struct {
	HouseNo     *string `json:"houseNo,omitempty" bson:"houseNo"`
	Subdistrict *string `json:"subdistrict,omitempty" bson:"subdistrict"`
	District    *string `json:"district,omitempty" bson:"district"`
	Province    *string `json:"province,omitempty" bson:"province"`
	PostalCode  *string `json:"postalCode,omitempty" bson:"postalCode"`
}

type EmployeeFilter struct {
	RowID     primitive.ObjectID `json:"rowId,omitempty" bson:"_id,omitempty"`
	IDNo      string             `json:"idNo,omitempty" bson:"idNo,omitempty"`
	FirstName string             `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName  string             `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Age       int                `json:"age,omitempty" bson:"age,omitempty"`
	Birthday  primitive.DateTime `json:"birthday,omitempty" bson:"birthday,omitempty"`
}

func GetCollectionEmployee(mongoPool *mongo.Client) *mongo.Collection {
	return mongoPool.Database(databaseName).Collection("employee")
}

var mongoErrCode errorcode.CrmErrorCode

func initMongoDB() (*logging.PatternLogger, DBPool) {
	appLogger := logging.InitInboundLogger("crm-util-go", logging.CrmDatabase)
	appLogger.Level = logging.LEVEL_ALL

	// Init Error Code
	errorcode.InitConfig("../config")
	mongoErrCode.SystemCode = "CIB"
	mongoErrCode.ModuleCode = "AS"

	var dbPool DBPool
	dbPool.DataSourceName = "mongodb://172.16.2.157:27017"
	// dbPool.DataSourceName = "mongodb://uat02appc:uat02appc1234@midmgdv1:29117,midmgdv2:29117,midmgdv3:29117/?authSource=mobileiddb-uat02&replicaSet=miduat&ssl=false"
	dbPool.MaxOpenConns = 100
	dbPool.MaxIdleConns = 5
	dbPool.MaxLifetime = time.Duration(3) * time.Minute
	dbPool.Logger = appLogger
	dbPool.AppName = "CRMApp"

	return appLogger, dbPool
}

func testCreateCollection() {
	transID := common.NewUUID()

	appLogger, dbPool := initMongoDB()
	startDT := appLogger.LogRequestDBClient(transID)

	_, err := dbPool.GetMongoDBPool(transID)

	if err != nil {
		appLogger.Error(transID, "Error Get MongoDB Pool")
		return
	}

	collName := "coll" + strconv.FormatInt(time.Now().Unix(), 10)
	appLogger.Info(transID, "collName", collName)

	singleResult := CreateCollection("test", collName)

	var result bson.M
	err = singleResult.Decode(&result)

	if err != nil {
		appLogger.Error(transID, "Error Run Command "+err.Error())
	} else {
		appLogger.Info(transID, "Success Run Command:", result)
	}

	appLogger.LogResponseDBClient(transID, "0", startDT)
}

func TestInsertEmployeeMongo(t *testing.T) {
	employeeResp := insertEmployee()

	if employeeResp.Code != "0" {
		t.Errorf("expecting result to be Code=0 but %s", employeeResp.Code)
	}
}

func insertEmployee() EmployeeResp {
	transID := common.NewUUID()

	var employeeResp EmployeeResp
	employeeResp.TransID = transID

	appLogger, dbPool := initMongoDB()
	startDT := appLogger.LogRequestDBClient(transID)

	mongoPool, err := dbPool.GetMongoDBPool(transID)

	if err != nil {
		appLogger.Error(transID, "Error Get MongoDB Pool")

		errorCodeResp := mongoErrCode.GenerateCRMDatabaseError(err.Error())

		employeeResp.Code = errorCodeResp.ErrorCode
		employeeResp.Msg = errorCodeResp.ErrorMessage

		return employeeResp
	}

	empColl := GetCollectionEmployee(mongoPool)

	dataSuffix := "100"

	empRecord := new(Employee)
	empRecord.IDNo = pointer.NewString(dataSuffix)
	empRecord.FirstName = pointer.NewString("Paravit " + dataSuffix)
	empRecord.LastName = pointer.NewString("Tunvichian " + dataSuffix)
	age, _ := strconv.ParseInt(dataSuffix, 10, 0)
	empRecord.Age = pointer.NewInt(int(age))
	empRecord.DepartmentID = StringToObjectIDPointer("6194b392cd0dcd12709f0acc")
	empRecord.Address = &Address{
		HouseNo:     pointer.NewString("HouseNo " + dataSuffix),
		Subdistrict: pointer.NewString("Subdistrict " + dataSuffix),
		District:    pointer.NewString("District " + dataSuffix),
		Province:    pointer.NewString("Province " + dataSuffix),
		PostalCode:  pointer.NewString("PostalCode " + dataSuffix),
	}

	birthdayTime, _ := timeUtil.StringToTime("yyyy-mm-ddThh:mi:ssZ", "2021-10-04T23:59:59Z")
	empRecord.Birthday = NewDateTimeMongo(birthdayTime)

	ctxInsert, cancelInsert := context.WithTimeout(context.Background(), mongoInsertTimeout)
	defer cancelInsert()

	insertResult, err := empColl.InsertOne(ctxInsert, empRecord)

	if err != nil {
		appLogger.Error(transID, "Error Insert "+err.Error())

		errorCodeResp := mongoErrCode.GenerateCRMDatabaseError("ERROR INSERT " + err.Error())

		employeeResp.Code = errorCodeResp.ErrorCode
		employeeResp.Msg = errorCodeResp.ErrorMessage

		return employeeResp
	} else {
		if insertResult != nil {
			objId, _ := insertResult.InsertedID.(primitive.ObjectID)

			empRecord.RowID = &objId
			appLogger.Info(transID, "Field _id=", *empRecord.RowID)
		}

		appLogger.Info(transID, "Insert Success", empRecord)
	}

	employeeResp.Code = "0"
	employeeResp.Msg = "Success"
	appLogger.LogResponseDBClient(transID, employeeResp.Code, startDT)

	return employeeResp
}

func TestUpdateEmployeeMongo(t *testing.T) {
	employeeResp := updateEmployee()

	if employeeResp.Code != "0" {
		t.Errorf("expecting result to be Code=0 but %s", employeeResp.Code)
	}
}

func updateEmployee() EmployeeResp {
	transID := common.NewUUID()

	var employeeResp EmployeeResp
	employeeResp.TransID = transID

	appLogger, dbPool := initMongoDB()
	startDT := appLogger.LogRequestDBClient(transID)

	empObjId, err := primitive.ObjectIDFromHex("61b1d27441745e392e3e2092")

	if err != nil {
		appLogger.Error(transID, "Error ObjectIDFromHex "+err.Error())

		errorCodeResp := mongoErrCode.GenerateAppError("crm-util-go", err.Error())

		employeeResp.Code = errorCodeResp.ErrorCode
		employeeResp.Msg = errorCodeResp.ErrorMessage

		return employeeResp
	}

	empFilter := new(EmployeeFilter)
	empFilter.RowID = empObjId
	//empFilter := bson.D{{"_id", empObjId}}

	mongoPool, err := dbPool.GetMongoDBPool(transID)

	if err != nil {
		appLogger.Error(transID, "Error Get MongoDB Pool")

		errorCodeResp := mongoErrCode.GenerateCRMDatabaseError(err.Error())

		employeeResp.Code = errorCodeResp.ErrorCode
		employeeResp.Msg = errorCodeResp.ErrorMessage

		return employeeResp
	}

	empColl := GetCollectionEmployee(mongoPool)

	empRecord := new(Employee)
	empRecord.IDNo = pointer.NewString("1")
	empRecord.FirstName = pointer.NewString("Paravit1")
	empRecord.LastName = pointer.NewString("Tunvichian1")
	empRecord.Age = pointer.NewInt(1)
	empRecord.Birthday = NewDateTimeMongo(time.Now())
	empRecord.DepartmentID = StringToObjectIDPointer("6194b2d2cd0dcd12709f0acb")

	updateBSON := bson.M{
		"$set": empRecord,
	}

	ctxUpdate, cancelUpdate := context.WithTimeout(context.Background(), mongoUpdateTimeout)
	defer cancelUpdate()

	optsUpsert := options.Update()
	optsUpsert.SetUpsert(false)

	updateResult, err := empColl.UpdateOne(ctxUpdate, empFilter, updateBSON, optsUpsert)

	if err != nil {
		appLogger.Error(transID, "Error Update "+err.Error())
		return employeeResp
	} else {
		appLogger.Info(transID, "Update Success", empRecord)
	}

	if updateResult != nil {
		appLogger.Debug(transID, "updateResult:", updateResult)

		if updateResult.MatchedCount == 0 {
			errorCodeResp := mongoErrCode.GenerateDataNotFound("Employee", "CRM")

			employeeResp.Code = errorCodeResp.ErrorCode
			employeeResp.Msg = errorCodeResp.ErrorMessage

			return employeeResp
		}
	}

	employeeResp.Code = "0"
	employeeResp.Msg = "Success"
	appLogger.LogResponseDBClient(transID, employeeResp.Code, startDT)

	return employeeResp
}

func TestQueryEmployeeMongo(t *testing.T) {
	employeeResp := queryEmployee()

	if employeeResp.Code != "0" {
		t.Errorf("expecting result to be Code=0 but %s", employeeResp.Code)
	}
}

func queryEmployee() EmployeeResp {
	transID := common.NewUUID()

	var employeeResp EmployeeResp
	employeeResp.TransID = transID

	appLogger, dbPool := initMongoDB()
	startDT := appLogger.LogRequestDBClient(transID)

	mongoPool, err := dbPool.GetMongoDBPool(transID)

	if err != nil {
		appLogger.Error(transID, "Error Get MongoDB Pool")

		errorCodeResp := mongoErrCode.GenerateCRMDatabaseError(err.Error())

		employeeResp.Code = errorCodeResp.ErrorCode
		employeeResp.Msg = errorCodeResp.ErrorMessage

		return employeeResp
	}

	empColl := GetCollectionEmployee(mongoPool)

	birthdayTime, err := timeUtil.StringToTime("yyyy-mm-ddThh:mi:ssZ07:00", "2020-10-01T00:00:00+07:00")

	if err != nil {
		appLogger.Error(transID, "Error StringToTime "+err.Error())
		return employeeResp
	}

	empFilter := make(map[string]interface{})
	objIdList := bson.A{}
	objIdList = append(objIdList, StringToObjectID("60c195ccdea661f99863eea4"))
	objIdList = append(objIdList, StringToObjectID("60c19392fc6d17518d710a5d"))
	objIdList = append(objIdList, StringToObjectID("60c19344601fd5b50bfe082d"))
	objIdList = append(objIdList, StringToObjectID("60c192d0239a62ac34a89561"))

	empFilter["_id"] = bson.D{{"$in", objIdList}}
	empFilter["birthday"] = bson.M{"$gte": NewDateTimeMongo(birthdayTime), "$lte": NewDateTimeMongo(time.Now())}

	findOneOptions := NewFindOneOptions(bson.D{{"age", -1}})

	ctxQuery, cancelQuery := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancelQuery()

	empResult := new(EmployeeDep)
	err = empColl.FindOne(ctxQuery, empFilter, findOneOptions).Decode(empResult)

	if err == nil {
		var empList []EmployeeDep
		empList = append(empList, *empResult)

		employeeResp.Code = "0"
		employeeResp.Msg = "Success"
		employeeResp.EmployeeList = empList
		employeeResp.TotalRecords = 1
	} else {
		var errorCodeResp errorcode.CrmErrorCodeResp
		if err == mongo.ErrNoDocuments {
			errorCodeResp = mongoErrCode.GenerateDataNotFound("Employee", "CRM")
		} else {
			errorCodeResp = mongoErrCode.GenerateCRMDatabaseError(err.Error())
		}

		employeeResp.Code = errorCodeResp.ErrorCode
		employeeResp.Msg = errorCodeResp.ErrorMessage
	}

	appLogger.LogResponseDBClient(transID, employeeResp.Code, startDT)
	appLogger.Info(transID, "Query Success. EmployeeList Size:", employeeResp.TotalRecords)
	appLogger.Info(transID, "EmployeeResp:", employeeResp)

	return employeeResp
}

func TestQueryEmployeeListMongo(t *testing.T) {
	employeeResp := queryEmployeeList()

	if employeeResp.Code != "0" {
		t.Errorf("expecting result to be Code=0 but %s", employeeResp.Code)
	}
}

func queryEmployeeList() EmployeeResp {
	transID := common.NewUUID()

	var employeeResp EmployeeResp
	employeeResp.TransID = transID

	appLogger, dbPool := initMongoDB()
	defer dbPool.Disconnect(transID)

	startDT := appLogger.LogRequestDBClient(transID)

	mongoPool, err := dbPool.GetMongoDBPool(transID)

	if err != nil {
		appLogger.Error(transID, "Error Get MongoDB Pool")

		errorCodeResp := mongoErrCode.GenerateCRMDatabaseError(err.Error())

		employeeResp.Code = errorCodeResp.ErrorCode
		employeeResp.Msg = errorCodeResp.ErrorMessage

		return employeeResp
	}

	empColl := GetCollectionEmployee(mongoPool)

	// Ex. birthday = UTC
	// birthdayTime, _ := common.StringToTime("yyyy-mm-ddThh:mi:ssZ", "2020-12-31T23:00:00Z")
	// Ex. birthday = UTC +07:00
	birthdayTime, err := timeUtil.StringToTime("yyyy-mm-ddThh:mi:ssZ07:00", "2021-10-01T00:00:00+07:00")

	if err != nil {
		appLogger.Error(transID, "Error StringToTime "+err.Error())
		return employeeResp
	}

	birthdayMongo := NewDateTimeMongo(birthdayTime)

	// Array
	objIdList := bson.A{}
	objIdList = append(objIdList, StringToObjectID("60c195ccdea661f99863eea4"))
	objIdList = append(objIdList, StringToObjectID("60c19392fc6d17518d710a5d"))
	objIdList = append(objIdList, StringToObjectID("60c19344601fd5b50bfe082d"))
	objIdList = append(objIdList, StringToObjectID("60c192d0239a62ac34a89561"))

	// Option 1: mapping struct
	//empFilter := new(EmployeeFilter)
	//empFilter.FirstName = "Paravit2"
	//empFilter.Birthday = *birthdayMongo

	// Option 2: mapping BSON  condition _id in ('xxx','yyy') and birthday >= xxx
	//empFilter := bson.M{
	//	"_id": bson.D{{"$in",objIdList}},
	//	"birthday": bson.M{"$gte": birthdayMongo},
	//}

	// Option 3: map bson.M
	//empFilter := make(map[string]interface{})
	//empFilter["_id"] = bson.D{{"$in",objIdList}}
	//empFilter["birthday"] = bson.M{"$gte": birthdayMongo, "$lte": db.NewDateTimeMongo(time.Now())}

	// Option 4: bson.D
	empFilter := bson.D{}
	empFilter = append(empFilter, bson.E{"_id", bson.D{{"$in", objIdList}}})
	empFilter = append(empFilter, bson.E{"birthday", bson.M{"$gte": birthdayMongo}})

	var limitRecords int64 = 2
	var pageNo int64 = 0
	findOptions := NewFindOptions(limitRecords, pageNo)

	// Option 1. sort := bson.D{{"birthday", -1}}
	// Option 2.
	ascending := 1
	descending := -1
	sort := bson.D{}
	sort = append(sort, bson.E{"birthday", descending})
	sort = append(sort, bson.E{"age", ascending})

	findOptions.SetSort(sort)

	ctxQuery, cancelQuery := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancelQuery()

	cursor, cursorErr := empColl.Find(ctxQuery, empFilter, findOptions)

	if cursorErr != nil {
		appLogger.Error(transID, "Error Cursor "+cursorErr.Error())
		return employeeResp
	}

	defer cursor.Close(context.TODO())

	var empList []EmployeeDep
	errDecode := cursor.All(context.TODO(), &empList)

	if errDecode != nil {
		appLogger.Error(transID, "Error Bind Object "+errDecode.Error())

		errorCodeResp := mongoErrCode.GenerateCRMDatabaseError(errDecode.Error())

		employeeResp.Code = errorCodeResp.ErrorCode
		employeeResp.Msg = errorCodeResp.ErrorMessage

		return employeeResp
	}

	employeeResp.Code = "0"
	employeeResp.Msg = "Success"
	employeeResp.EmployeeList = empList
	employeeResp.TotalRecords = int64(len(empList))

	appLogger.LogResponseDBClient(transID, employeeResp.Code, startDT)

	appLogger.Info(transID, "Query Success. EmployeeList Size:", employeeResp.TotalRecords)
	appLogger.Info(transID, "EmployeeResp:", employeeResp)

	return employeeResp
}

func TestQueryAggregateListMongo(t *testing.T) {
	employeeResp := queryAggregateListMongo()

	if employeeResp.Code != "0" {
		t.Errorf("expecting result to be Code=0 but %s", employeeResp.Code)
	}
}

func queryAggregateListMongo() EmployeeResp {
	transID := common.NewUUID()

	var employeeResp EmployeeResp
	employeeResp.TransID = transID

	appLogger, dbPool := initMongoDB()
	startDT := appLogger.LogRequestDBClient(transID)

	mongoPool, err := dbPool.GetMongoDBPool(transID)

	if err != nil {
		appLogger.Error(transID, "Error Get MongoDB Pool")

		errorCodeResp := mongoErrCode.GenerateCRMDatabaseError(err.Error())

		employeeResp.Code = errorCodeResp.ErrorCode
		employeeResp.Msg = errorCodeResp.ErrorMessage

		return employeeResp
	}

	empColl := GetCollectionEmployee(mongoPool)

	ctxQuery, cancelQuery := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancelQuery()

	/*
		$lookup : left join
		$match : filter
			{{"$match", bson.D{{"firstName", bson.D{{"$eq", "Paravit 1"}}}}}},
			{{"$match", bson.D{{"firstName", "Paravit 1"}}}}
			{{"$match", bson.D{{"department.departmentName", bson.D{{"$eq", "IT Dev"}}}}}},
	*/
	pipeLine := mongo.Pipeline{
		{{"$match", bson.D{{"firstName", "Paravit1"}}}},
		{{"$lookup", bson.D{
			{"from", "department"},
			{"localField", "departmentID"},
			{"foreignField", "_id"},
			{"as", "department"},
		}}},
	}

	opts := options.Aggregate().SetMaxTime(mongoQueryTimeout)

	cursor, cursorErr := empColl.Aggregate(ctxQuery, pipeLine, opts)

	if cursorErr != nil {
		appLogger.Error(transID, "Error Cursor "+cursorErr.Error())
		return employeeResp
	}

	defer cursor.Close(context.TODO())

	var empList []EmployeeDep
	errDecode := cursor.All(context.TODO(), &empList)

	if errDecode != nil {
		appLogger.Error(transID, "Error Bind Object "+errDecode.Error())
		return employeeResp
	}

	employeeResp.Code = "0"
	employeeResp.Msg = "Success"
	employeeResp.EmployeeList = empList
	employeeResp.TotalRecords = int64(len(empList))

	appLogger.LogResponseDBClient(transID, employeeResp.Code, startDT)

	appLogger.Info(transID, "Query Success. EmployeeList Size:", employeeResp.TotalRecords)
	appLogger.Info(transID, "EmployeeResp:", employeeResp)

	return employeeResp
}

func TestDeleteEmployeeMongo(t *testing.T) {
	employeeResp := deleteEmployee()

	if employeeResp.Code != "0" {
		t.Errorf("expecting result to be Code=0 but %s", employeeResp.Code)
	}
}

func deleteEmployee() EmployeeResp {
	transID := common.NewUUID()

	var employeeResp EmployeeResp
	employeeResp.TransID = transID

	appLogger, dbPool := initMongoDB()
	startDT := appLogger.LogRequestDBClient(transID)

	objId, err := primitive.ObjectIDFromHex("aaa194066d05ffac04f4261f")

	if err != nil {
		appLogger.Error(transID, "Error ObjectIDFromHex "+err.Error())

		errorCodeResp := mongoErrCode.GenerateAppError("crm-util-go", err.Error())

		employeeResp.Code = errorCodeResp.ErrorCode
		employeeResp.Msg = errorCodeResp.ErrorMessage

		return employeeResp
	}

	mongoPool, err := dbPool.GetMongoDBPool(transID)

	if err != nil {
		appLogger.Error(transID, "Error Get MongoDB Pool")

		errorCodeResp := mongoErrCode.GenerateCRMDatabaseError(err.Error())

		employeeResp.Code = errorCodeResp.ErrorCode
		employeeResp.Msg = errorCodeResp.ErrorMessage

		return employeeResp
	}

	empColl := GetCollectionEmployee(mongoPool)

	empFilter := new(EmployeeFilter)
	empFilter.RowID = objId

	ctxDelete, cancelDelete := context.WithTimeout(context.Background(), mongoDeleteTimeout)
	defer cancelDelete()

	deleteResult, err := empColl.DeleteOne(ctxDelete, empFilter)

	// DeleteMany Important: An empty document (e.g. bson.D{}) should be used to delete all documents in the collection.
	//deleteResult, err := empColl.DeleteMany(ctxDelete, bson.D{})

	if err != nil {
		appLogger.Error(transID, "ERROR DELETE "+err.Error())

		errorCodeResp := mongoErrCode.GenerateCRMDatabaseError(err.Error())

		employeeResp.Code = errorCodeResp.ErrorCode
		employeeResp.Msg = errorCodeResp.ErrorMessage

		return employeeResp
	}

	appLogger.Info(transID, "Deleted Document:", deleteResult.DeletedCount)

	employeeResp.Code = "0"
	employeeResp.Msg = "Success"
	employeeResp.TotalRecords = deleteResult.DeletedCount

	appLogger.LogResponseDBClient(transID, employeeResp.Code, startDT)
	appLogger.Info(transID, "EmployeeResp:", employeeResp)

	return employeeResp
}

func TestInsertTransactionsACIDMongo(t *testing.T) {
	employeeResp := insertTransactionsACID()

	if employeeResp.Code != "0" {
		t.Errorf("expecting result to be Code=0 but %s", employeeResp.Code)
	}
}

func insertTransactionsACID() EmployeeResp {
	// (IllegalOperation) Transaction numbers are only allowed on a replica set member or mongos
	// The mongos instances provide the interface between the client applications and the sharded cluster.

	transID := common.NewUUID()

	var employeeResp EmployeeResp
	employeeResp.TransID = transID

	appLogger, dbPool := initMongoDB()
	startDT := appLogger.LogRequestDBClient(transID)

	// Provide data
	empRecord1 := new(Employee)
	empRecord1.IDNo = pointer.NewString("100")
	empRecord1.FirstName = pointer.NewString("Paravit100")
	empRecord1.LastName = pointer.NewString("Tunvichian100")
	empRecord1.Age = pointer.NewInt(100)
	empRecord1.DepartmentID = StringToObjectIDPointer("6194b2d2cd0dcd12709f0acb")
	birthdayTime1, _ := timeUtil.StringToTime("yyyy-mm-ddThh:mi:ssZ", "2021-11-01T23:59:59Z")
	empRecord1.Birthday = NewDateTimeMongo(birthdayTime1)

	/*
		Rollback transaction because required idNo
		ERROR -> write exception: write errors: [Document failed validation])
	*/
	empRecord2 := new(Employee)
	empRecord2.IDNo = pointer.NewString("200")
	empRecord2.FirstName = pointer.NewString("Paravit200")
	empRecord2.LastName = pointer.NewString("Tunvichian200")
	empRecord2.Age = pointer.NewInt(200)
	empRecord2.DepartmentID = StringToObjectIDPointer("6194b2d2cd0dcd12709f0acb")
	birthdayTime2, _ := timeUtil.StringToTime("yyyy-mm-ddThh:mi:ssZ", "2021-11-02T23:59:59Z")
	empRecord2.Birthday = NewDateTimeMongo(birthdayTime2)

	mongoPool, err := dbPool.GetMongoDBPool(transID)

	if err != nil {
		appLogger.Error(transID, "Error Get MongoDB Pool "+err.Error())

		errorCodeResp := mongoErrCode.GenerateCRMDatabaseError(err.Error())

		employeeResp.Code = errorCodeResp.ErrorCode
		employeeResp.Msg = errorCodeResp.ErrorMessage

		return employeeResp
	}

	session, err := mongoPool.StartSession()

	if err != nil {
		appLogger.Error(transID, "Error Start New Session "+err.Error())

		errorCodeResp := mongoErrCode.GenerateCRMDatabaseError(err.Error())

		employeeResp.Code = errorCodeResp.ErrorCode
		employeeResp.Msg = errorCodeResp.ErrorMessage

		return employeeResp
	}

	defer session.EndSession(context.Background())

	callbackFunc := func(sessionContext mongo.SessionContext) error {
		transOpts := NewTransactionOptions()

		err = session.StartTransaction(transOpts)

		if err != nil {
			return err
		}

		empColl := GetCollectionEmployee(mongoPool)

		// Transaction 1
		result, err := empColl.InsertOne(sessionContext, empRecord1)

		if err != nil {
			return err
		}

		appLogger.Info(transID, "Transaction 1 -> Object id", result.InsertedID)

		// Transaction 2
		result, err = empColl.InsertOne(sessionContext, empRecord2)

		if err != nil {
			return err
		}

		err = session.CommitTransaction(sessionContext)

		if err != nil {
			appLogger.Error(transID, "Error CommitTransaction "+err.Error())
			return err
		}

		appLogger.Info(transID, "Transaction 2 -> Object id", result.InsertedID)
		return nil
	}

	err = mongo.WithSession(context.Background(), session, callbackFunc)

	if err != nil {
		appLogger.Error(transID, "Error "+err.Error())

		abortErr := session.AbortTransaction(context.Background())

		if abortErr != nil {
			appLogger.Error(transID, "Error AbortTransaction "+abortErr.Error())
			panic(abortErr)
		}

		errorCodeResp := mongoErrCode.GenerateCRMDatabaseError(err.Error())

		employeeResp.Code = errorCodeResp.ErrorCode
		employeeResp.Msg = errorCodeResp.ErrorMessage
	} else {
		employeeResp.Code = "0"
		employeeResp.Msg = "Success"
	}

	appLogger.LogResponseDBClient(transID, employeeResp.Code, startDT)
	appLogger.Info(transID, "EmployeeResp:", employeeResp)

	return employeeResp
}

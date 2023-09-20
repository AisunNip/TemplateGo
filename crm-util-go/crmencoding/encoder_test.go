package crmencoding_test

import (
	"crm-util-go/crmencoding"
	"crm-util-go/pointer"
	"fmt"
	"testing"
	"time"
)

func TestEncodeBase64(t *testing.T) {
	msg := "Test@12345$ทดสอบ"
	encodedMsg := crmencoding.EncodeBase64(msg)
	fmt.Println("EncodedMsg: " + encodedMsg)

	planText, encError := crmencoding.DecodeBase64(encodedMsg)
	if encError != nil {
		t.Errorf("DecodeBase64 Error %s", encError.Error())
	}
	fmt.Println("PlanText: " + planText)

	if planText != msg {
		t.Error("TestEncodeBase64 Error")
	}
}

func TestEncodeMultiQueryString(t *testing.T) {
	query := make(map[string]string)
	query["firstName"] = "Paravit"
	query["lastName"] = "Tunvichian"

	multiQueryString := crmencoding.EncodeMultiQueryString(query)
	fmt.Println("MultiQueryString: " + multiQueryString)

	if multiQueryString != "firstName=Paravit&lastName=Tunvichian" {
		t.Error("EncodeMultiQueryString Error")
	}
}

func TestEncodeURL(t *testing.T) {
	query := make(map[string]string)
	query["firstName"] = "Paravit"
	query["lastName"] = "Tunvichian"

	encodedURL, urlErr := crmencoding.EncodeURL("http://www.google.com", "abc/xx yy", query)

	if urlErr != nil {
		t.Errorf("EncodeURL Error %s", urlErr.Error())
	} else {
		fmt.Println("EncodeURL: " + encodedURL)

		if encodedURL != "http://www.google.com/abc/xx%20yy?firstName=Paravit&lastName=Tunvichian" {
			t.Error("EncodeURL Error")
		}
	}
}

func TestMapToJson(t *testing.T) {
	myMap := make(map[string]interface{})
	myMap["firstName"] = "Paravit"
	myMap["lastName"] = "Tunvichian"
	myMap["age"] = 40
	myMap["birthDate"] = time.Now()

	json, err := crmencoding.MapToJsonIndent(myMap)

	if err != nil {
		t.Errorf("MapToJsonIndent Error %s", err.Error())
	} else {
		fmt.Println("Json:", json)
	}
}

func TestJsonToMap(t *testing.T) {
	json := `{"age":40,"birthDate":"2021-07-04T23:59:59+07:00","firstName":"Paravit","lastName":"Tunvichian"}`
	mapResult, err := crmencoding.JsonToMap(json)

	if err != nil {
		t.Errorf("JsonToMap Error %s", err.Error())
	} else {
		fmt.Println("Map:", mapResult)
	}
}

type MyName struct {
	FirstName *string    `json:"firstName,omitempty"`
	LastName  *string    `json:"lastName,omitempty"`
	BirthDate *time.Time `json:"birthDate,omitempty"`
	Age       *int       `json:"age,omitempty"`
}

func TestJsonToStruct(t *testing.T) {
	json := `{"birthDate":"2021-07-04T23:59:59+07:00","firstName":"Paravit","lastName": null, "age": null}`
	myName := new(MyName)
	err := crmencoding.JsonToStruct(json, myName)

	if err != nil {
		t.Errorf("JsonToStruct Error %s", err.Error())
	} else {
		fmt.Println("firstName:", pointer.GetStringValue(myName.FirstName))
		fmt.Println("lastName:", pointer.GetStringValue(myName.LastName))
		fmt.Println("birthDate:", pointer.GetTimeValue(myName.BirthDate))
		fmt.Println("age:", pointer.GetIntValue(myName.Age))
	}
}

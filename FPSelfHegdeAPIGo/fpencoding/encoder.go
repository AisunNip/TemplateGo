package fpencoding

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/url"
)

func EncodeBase64(msg string) string {
	return base64.StdEncoding.EncodeToString([]byte(msg))
}

func DecodeBase64(encodedMsg string) (planText string, err error) {
	var decoded []byte
	decoded, err = base64.StdEncoding.DecodeString(encodedMsg)

	if err != nil {
		return
	}

	planText = string(decoded)

	return
}

func EncodeQueryString(query string) string {
	return url.QueryEscape(query)
}

func EncodeMultiQueryString(query map[string]string) string {
	var encodedQuery string

	if len(query) > 0 {
		params := url.Values{}

		for key, val := range query {
			params.Add(key, val)
		}

		encodedQuery = params.Encode()
	}

	return encodedQuery
}

func EncodeURLPath(path string) string {
	return url.PathEscape(path)
}

func EncodeURL(baseURL string, path string, query map[string]string) (string, error) {
	baseUrl, err := url.Parse(baseURL)

	if err != nil {
		err = errors.New("Malformed URL: " + err.Error())
		return "", err
	}

	schemaArray := [7]string{"http", "https", "ftp", "mailto", "file", "data", "irc"}

	isValidScheme := false
	for _, value := range schemaArray {
		if baseUrl.Scheme == value {
			isValidScheme = true
			break
		}
	}

	if !isValidScheme {
		err = errors.New("Malformed URL is invalid scheme")
		return "", err
	}

	baseUrl.Path += path
	baseUrl.RawQuery = EncodeMultiQueryString(query)

	return baseUrl.String(), err
}

func EncodeFormBody(body map[string]string) string {
	var encodedData string

	if body != nil {
		formData := url.Values{}

		for k, v := range body {
			formData.Add(k, v)
		}

		encodedData = formData.Encode()
	}

	return encodedData
}

func structToJson(structInput interface{}, isIndent bool) (jsonString string, err error) {
	var binary []byte

	if isIndent {
		binary, err = json.MarshalIndent(structInput, "", "  ")
	} else {
		binary, err = json.Marshal(structInput)
	}

	if err != nil {
		return
	}

	jsonString = string(binary)
	return
}

func StructToJson(structInput interface{}) (jsonString string, err error) {
	return structToJson(structInput, false)
}

func StructToJsonIndent(structInput interface{}) (jsonString string, err error) {
	return structToJson(structInput, true)
}

func JsonToStruct(jsonString string, structOutput interface{}) error {
	return json.Unmarshal([]byte(jsonString), structOutput)
}

func MapToJson(dataMap map[string]interface{}) (jsonString string, err error) {
	return structToJson(dataMap, false)
}

func MapToJsonIndent(dataMap map[string]interface{}) (jsonString string, err error) {
	return structToJson(dataMap, true)
}

func JsonToMap(jsonString string) (map[string]interface{}, error) {
	var dataMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &dataMap)
	return dataMap, err
}
func XmlToStruct(xmlString string, structOutput interface{}) error {
	return xml.Unmarshal([]byte(xmlString), structOutput)
}

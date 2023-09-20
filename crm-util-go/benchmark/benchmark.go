package benchmark

import (
	"bytes"
	"encoding/json"
	"io"
)

func BMIOReadAll(reader io.Reader) (map[string]interface{}, error) {
	var (
		m    map[string]interface{}
		b, _ = io.ReadAll(reader)
	)

	return m, json.Unmarshal(b, &m)
}

func BMIOCopy(reader io.Reader) (map[string]interface{}, error) {
	var (
		m    map[string]interface{}
		buf  bytes.Buffer
		_, _ = io.Copy(&buf, reader)
	)

	return m, json.Unmarshal(buf.Bytes(), &m)
}

func BMHasStringValueLen(data string) bool {
	if len(data) > 0 {
		return true
	} else {
		return false
	}
}

func BMHasStringValueEqual(data string) bool {
	if data == "" {
		return false
	} else {
		return true
	}
}

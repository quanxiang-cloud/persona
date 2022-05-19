package utils

import (
	"bytes"
	"encoding/json"
)

// Struct2Bytes 结构体转换为字节
func Struct2Bytes(reqData interface{}) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqData)
	if err != nil {
		return nil, err
	}
	return &buf, err
}

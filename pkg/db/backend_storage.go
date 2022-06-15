package db

import (
	"context"
	"encoding/json"
)

// BackendStorage 抽象接口
type BackendStorage interface {
	Put(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (map[string]string, error)
	GetWithPrefix(ctx context.Context, key string) ([]ImportReqData, error)
	PutWithVersion(ctx context.Context, version string, key string, value string) error
	GetWithVersion(ctx context.Context, version string, key string) (map[string]string, error)
	UserPutWithVersion(ctx context.Context, version string, key string, value string) error
	UserGetWithVersion(ctx context.Context, version string, key string) (map[string]string, error)
	PutData(ctx *context.Context, key *string, value interface{}) error
	GetData(ctx *context.Context, key *string) (*json.RawMessage, error)
	UpdateData(ctx *context.Context, key *string, value interface{}) error
	GetDataByKVs(ctx *context.Context, kvs *map[string]interface{}) ([]*json.RawMessage, error)
	DeleteData(ctx *context.Context, key *string) error
	SearchWithKey(ctx context.Context, key string) (interface{}, error)
}

// Kv Kv
type Kv struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Version  string `json:"version"` // 使用方维护
	UserID   string `json:"user_id"`
	DataType string `json:"data_type"` // 内部使用
}

// ImportReqData 导入数据请求
type ImportReqData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

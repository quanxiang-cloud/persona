package model

import (
	"fmt"
	"git.internal.yunify.com/qxp/persona/pkg/config"
	pes "git.internal.yunify.com/qxp/persona/pkg/db/elasticsearch"
	petcd "git.internal.yunify.com/qxp/persona/pkg/db/etcd"
)

// DBFactory 根据配置不同返回不同的db对象
func DBFactory(conf *config.Configs) (interface{}, error) {
	switch conf.BackendStorage {
	case "es":
		b, err := pes.NewEs(conf)
		if err != nil {
			return nil, err
		}
		return b, nil
	case "etcd":
		b, err := petcd.NewEtcd(conf)
		if err != nil {
			return nil, err
		}
		return b, nil
	default:
		panic(fmt.Sprintf("Unsupported backend of: %s", conf.BackendStorage))
	}
}

// DataSet DataSet
type DataSet struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	Type      int64  `json:"type"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	DataType  string `json:"data_type"`
}

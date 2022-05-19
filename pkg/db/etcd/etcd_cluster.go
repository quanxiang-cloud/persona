package db

import (
	"context"
	"encoding/json"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/persona/pkg/config"
	"git.internal.yunify.com/qxp/persona/pkg/db"

	"go.etcd.io/etcd/clientv3"
	"strings"
	"time"
)

// Etcd etcd
type Etcd struct {
	client     *clientv3.Client
	etcdConfig config.EtcdConfig
	prefix     string
}

// DeleteData not need implement
func (d *Etcd) DeleteData(ctx *context.Context, key *string) error {
	return nil
}

// UpdateData not need implement
func (d *Etcd) UpdateData(ctx *context.Context, key *string, value interface{}) error {
	return nil
}

// GetDataByKVs not need implement
func (d *Etcd) GetDataByKVs(ctx *context.Context, kvs *map[string]interface{}) ([]*json.RawMessage, error) {
	return nil, nil
}

// PutData not need implement
func (d *Etcd) PutData(ctx *context.Context, key *string, value interface{}) error {
	return nil
}

// GetData not need implement
func (d *Etcd) GetData(ctx *context.Context, key *string) (*json.RawMessage, error) {
	return nil, nil
}

// NewEtcdClient new etcd client
func NewEtcdClient(config config.EtcdConfig) (*clientv3.Client, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Addrs,
		DialTimeout: config.Timeout * time.Second,
		Username:    config.Username,
		Password:    config.Password,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

// NewEtcd new etcd
func NewEtcd(conf *config.Configs) (db.BackendStorage, error) {
	cli, err := NewEtcdClient(conf.Etcd)
	return &Etcd{
		client:     cli,
		etcdConfig: conf.Etcd,
		prefix:     conf.HostName,
	}, err
}

// Put 存数据
func (d *Etcd) Put(ctx context.Context, key string, value string) error {
	_, err := d.client.Put(ctx, d.addPrefix(key), value)
	return err
}

// Get 取数据
func (d *Etcd) Get(ctx context.Context, key string) (map[string]string, error) {
	res, err := d.client.Get(ctx, d.addPrefix(key))
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, 0)
	for _, ev := range res.Kvs {
		result[d.removePrefix(string(ev.Key))] = string(ev.Value)
	}
	return result, nil
}

// GetWithPrefix 获取前缀列表
func (d *Etcd) GetWithPrefix(ctx context.Context, key string) ([]db.ImportReqData, error) {
	res, err := d.client.Get(ctx, d.addPrefix(key), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	result := make([]db.ImportReqData, 0)
	for _, ev := range res.Kvs {
		result = append(result, db.ImportReqData{
			Key:   d.removePrefix(string(ev.Key)),
			Value: string(ev.Value),
		})
	}
	return result, nil
}

// PutWithVersion 存储带前缀的key
func (d *Etcd) PutWithVersion(ctx context.Context, version string, key string, value string) error {
	_, err := d.client.Put(ctx, d.addPrefix2New(version, key), value)
	return err
}

// GetWithVersion 获取带版本的value
func (d *Etcd) GetWithVersion(ctx context.Context, version string, key string) (map[string]string, error) {
	res, err := d.client.Get(ctx, d.addPrefix2New(version, key))
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, 0)
	if len(res.Kvs) > 0 {
		for _, ev := range res.Kvs {
			result[d.removePrefix2New(version, string(ev.Key))] = string(ev.Value)
		}
		return result, nil
	}

	// Compatible with old formats
	res, err = d.client.Get(ctx, d.addPrefix2(version, key))
	if err != nil {
		return nil, err
	}
	if len(res.Kvs) > 0 {
		for _, ev := range res.Kvs {
			result[d.removePrefix2(version, string(ev.Key))] = string(ev.Value)
		}
	}

	return result, nil
}

// UserPutWithVersion 存储用户版本
func (d *Etcd) UserPutWithVersion(ctx context.Context, version string, key string, value string) error {
	userID := logger.STDHeader(ctx)["User-Id"]
	k := d.addPrefix3(userID, version, key)
	_, err := d.client.Put(ctx, k, value)
	return err
}

// UserGetWithVersion 获取用户版本
func (d *Etcd) UserGetWithVersion(ctx context.Context, version string, key string) (map[string]string, error) {
	userID := logger.STDHeader(ctx)["User-Id"]
	k := d.addPrefix3(userID, version, key)
	res, err := d.client.Get(ctx, k)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, 0)
	for _, ev := range res.Kvs {
		result[d.removePrefix3(userID, version, string(ev.Key))] = string(ev.Value)
	}
	return result, nil
}

// addPrefix 添加前缀
func (d *Etcd) addPrefix(key string) string {
	return d.prefix + "_" + key
}

// removePrefix 移除前缀
func (d *Etcd) removePrefix(key string) string {
	pre := d.prefix + "_"
	if strings.HasPrefix(key, pre) {
		return key[len(pre):]
	}
	return key
}

func (d *Etcd) addPrefix2New(version string, key string) string {
	return d.prefix + "_" + key + "_" + version
}

func (d *Etcd) removePrefix2New(version string, key string) string {
	return key[len(d.prefix)+1 : len(key)-len(version)-1]
}

func (d *Etcd) addPrefix2(version string, key string) string {
	return d.prefix + "_" + version + "_" + key
}

func (d *Etcd) removePrefix2(version string, key string) string {
	pre := d.prefix + "_" + version + "_"
	if strings.HasPrefix(key, pre) {
		return key[len(pre):]
	}
	return key
}

func (d *Etcd) addPrefix3(userID string, version string, key string) string {
	return d.prefix + "_" + userID + "_" + version + "_" + key
}

func (d *Etcd) removePrefix3(userID string, version string, key string) string {
	pre := d.prefix + "_" + userID + "_" + version + "_"
	if strings.HasPrefix(key, pre) {
		return key[len(pre):]
	}
	return key
}

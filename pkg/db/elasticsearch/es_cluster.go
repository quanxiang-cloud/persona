package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"git.internal.yunify.com/qxp/persona/pkg/config"
	"git.internal.yunify.com/qxp/persona/pkg/db"
	"git.internal.yunify.com/qxp/persona/pkg/misc/elastic2"
	"git.internal.yunify.com/qxp/persona/pkg/misc/logger"
	"github.com/olivere/elastic/v7"
)

var (
	// EsClient EsClient
	EsClient *elastic.Client
	// MaxPageSize es单次最大返回条数
	MaxPageSize = 1000
	// TypeOfDataSet es中的数据集类型
	TypeOfDataSet = "dataSet"
	// TypeOfDefault es中的默认数据类型
	TypeOfDefault = "default"
)

// NewClient new elasticsearch client
func NewClient(conf *elastic2.Config, opts ...elastic.ClientOptionFunc) (*elastic.Client, error) {
	for _, host := range conf.Host {
		opts = append(opts, elastic.SetURL(host))
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	_, _, err = client.Ping(conf.Host[0]).Do(context.Background())
	if err != nil {
		return nil, err
	}

	_, err = client.ElasticsearchVersion(conf.Host[0])
	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewEsClient return a new es client
func NewEsClient(config *config.ESConf) (*elastic.Client, error) {
	if EsClient != nil {
		return EsClient, nil
	}
	conf := elastic2.Config{
		Host: config.Host,
		Log:  false,
	}
	client, err := NewClient(&conf, elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	EsClient = client

	return client, nil
}

// Elasticsearch es
type Elasticsearch struct {
	client   *elastic.Client
	esConfig *config.ESConf
	prefix   string
}

// NewEs new es
func NewEs(conf *config.Configs) (db.BackendStorage, error) {
	cli, err := NewEsClient(&conf.ES)

	return &Elasticsearch{
		client:   cli,
		esConfig: &conf.ES,
		prefix:   conf.HostName,
	}, err
}

// Put 存储v到key
func (d *Elasticsearch) Put(ctx context.Context, key string, value string) error {
	data := db.Kv{
		Key:      key,
		Value:    value,
		DataType: TypeOfDefault,
	}
	return d.PutData(&ctx, &key, &data)
}

// Get 获取key的值
// 最终返回json
func (d *Elasticsearch) Get(ctx context.Context, key string) (map[string]string, error) {
	res, err := d.GetData(&ctx, &key)
	var resp = make(map[string]string)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return resp, nil
	}
	if err := json.Unmarshal(*res, &resp); err != nil {
		return nil, err
	}

	return resp, err
}

// PutWithVersion 带版本的数据
func (d *Elasticsearch) PutWithVersion(ctx context.Context, version string, key string, value string) error {
	k := d.genIDAndVersion(&key, &version)
	data := db.Kv{
		Key:      k,
		Value:    value,
		Version:  version,
		DataType: TypeOfDefault,
	}
	return d.PutData(&ctx, &k, &data)
}

// GetWithPrefix 返回匹配前缀的数据
// 默认情况下只返回1000条。如果大于这个数，则循环去es取
func (d *Elasticsearch) GetWithPrefix(ctx context.Context, key string) ([]db.ImportReqData, error) {
	var q = map[string]string{
		"key": key,
	}
	var Offset = 0
	var FoundTotal = 0
	var result = make([]db.ImportReqData, 0)
SEARCH:
	search := d.client.Search().Index(d.esConfig.DefaultIndex)
	d.PrefixQuery(search, &q)
	ret, err := search.From(Offset).Size(MaxPageSize).Do(ctx)
	if err != nil {
		return nil, err
	}
	for _, r := range ret.Hits.Hits {
		var res db.ImportReqData
		err := json.Unmarshal(r.Source, &res)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	FoundTotal += MaxPageSize
	if (ret.Hits.TotalHits.Value - int64(FoundTotal)) > 0 {
		Offset += MaxPageSize
		goto SEARCH
	}

	return result, nil
}

// GetWithVersion 获取带版本数据
func (d *Elasticsearch) GetWithVersion(ctx context.Context, version string, key string) (map[string]string, error) {
	id := d.genIDAndVersion(&key, &version)
	resp := make(map[string]string)
	r, err := d.Get(ctx, id)
	if r != nil {
		resp[key] = r["value"]
	}
	return resp, err
}

// UserPutWithVersion 设置用户带版本的值
func (d *Elasticsearch) UserPutWithVersion(ctx context.Context, version string, key string, value string) error {
	userID := logger.STDHeader(ctx)["User-Id"]
	k := d.genIDVersionAndUserID(&key, &version, &userID)
	data := db.Kv{
		Key:      k,
		Value:    value,
		Version:  version,
		UserID:   userID,
		DataType: TypeOfDefault,
	}
	return d.PutData(&ctx, &k, &data)
}

// UserGetWithVersion 获取用户带版本的值
func (d *Elasticsearch) UserGetWithVersion(ctx context.Context, version string, key string) (map[string]string, error) {
	userID := logger.STDHeader(ctx)["User-Id"]
	k := d.genIDVersionAndUserID(&key, &version, &userID)
	r, err := d.Get(ctx, k)
	resp := make(map[string]string)
	if len(r) > 0 {
		resp[key] = r[k]
	}

	return resp, err
}

// genIDAndVersion 生成es中需要的ID
// format is: {key}_{version}
func (d *Elasticsearch) genIDAndVersion(key *string, version *string) string {
	return fmt.Sprintf("%s_%s", *key, *version)
}

// genIDVersionAndUserID 生成es中某个用户对应的key及version
func (d *Elasticsearch) genIDVersionAndUserID(key *string, version *string, UserID *string) string {
	return fmt.Sprintf("%s_%s_%s", *UserID, *version, *key)
}

// CreateIndex 初始化索引
func (d *Elasticsearch) CreateIndex(ctx context.Context, index string) error {
	if len(index) == 0 {
		return errors.New("index can not be none")
	}
	exists, err := d.CheckIndexExists(ctx, index)
	if err != nil {
		return err
	}
	if exists == true && err == nil {
		return fmt.Errorf("index %s is already exists", index)
	}
	d.client.CreateIndex(index)

	return nil
}

// CheckIndexExists 检查index是否存在
func (d *Elasticsearch) CheckIndexExists(ctx context.Context, index string) (bool, error) {
	exists, err := d.client.IndexExists(index).Do(ctx)
	if err != nil {
		return false, nil
	}
	if !exists {
		return true, nil
	}
	return false, nil
}

// AndQueryCondition es and查询过滤条件.
// conditions: {"k": "v"}
func (d *Elasticsearch) AndQueryCondition(Query *elastic.SearchService, conditions *map[string]interface{}) *elastic.SearchService {
	if conditions == nil {
		Query = Query.Query(elastic.NewMatchAllQuery())
		return Query
	}
	q := elastic.NewBoolQuery()
	for k, v := range *conditions {
		q = q.Must(elastic.NewTermQuery(k, v))
	}
	Query = Query.Query(q)
	return Query
}

// PrefixQuery 前缀查询
// 返回所有前缀匹配的value
func (d *Elasticsearch) PrefixQuery(Query *elastic.SearchService, conditions *map[string]string) *elastic.SearchService {
	var q *elastic.PrefixQuery
	for k, v := range *conditions {
		q = elastic.NewPrefixQuery(k, v)
	}
	return Query.Query(q)
}

// InitEsIndex 初始化index
func InitEsIndex(conf *config.Configs) error {
	client, err := NewEsClient(&conf.ES)
	if err != nil {
		return err
	}
	ctx := context.Background()
	exists, err := client.IndexExists(conf.ES.DefaultIndex).Do(ctx)
	// 没有默认index则创建
	if !exists {
		_, err := client.CreateIndex(conf.ES.DefaultIndex).BodyString(IndexMappingLatest).Do(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// PutData 存储v到key
func (d *Elasticsearch) PutData(ctx *context.Context, key *string, value interface{}) error {
	_, err := d.client.
		Index().
		Index(d.esConfig.DefaultIndex).
		Id(*key).
		BodyJson(value).
		Do(*ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetData 获取key的值
func (d *Elasticsearch) GetData(ctx *context.Context, key *string) (*json.RawMessage, error) {
	res, err := d.client.
		Get().
		Index(d.esConfig.DefaultIndex).
		Id(*key).
		Do(*ctx)
	if err != nil {
		return nil, err
	}
	if res.Found == false {
		return &res.Source, nil
	}
	return &res.Source, nil
}

// UpdateData 更新数据
// value 可传map或struct
func (d *Elasticsearch) UpdateData(ctx *context.Context, key *string, value interface{}) error {
	_, err := d.client.
		Update().
		Index(d.esConfig.DefaultIndex).
		Id(*key).
		Doc(value).
		Do(*ctx)
	if err != nil {
		return err
	}
	return nil
}

// DeleteData 根据key删除数据
func (d *Elasticsearch) DeleteData(ctx *context.Context, key *string) error {
	_, err := d.client.
		Delete().
		Index(d.esConfig.DefaultIndex).
		Id(*key).
		Refresh("true").
		Do(*ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetDataByKVs 根据k v过滤数据。一次返回所有数据
func (d *Elasticsearch) GetDataByKVs(ctx *context.Context, kvs *map[string]interface{}) ([]*json.RawMessage, error) {
	if kvs == nil {
		return nil, errors.New("GetDataByKVs: need one or more condition(s)")
	}
	var Offset = 0
	var FoundTotal = 0
	var Resp = make([]*json.RawMessage, 0)
SEARCH:
	search := d.client.Search().Index(d.esConfig.DefaultIndex)
	search = d.AndQueryCondition(search, kvs)
	ret, err := search.From(Offset).Size(MaxPageSize).Do(*ctx)
	if err != nil {
		return nil, err
	}
	for _, r := range ret.Hits.Hits {
		Resp = append(Resp, &(r.Source))
	}
	FoundTotal += MaxPageSize
	if (ret.Hits.TotalHits.Value - int64(FoundTotal)) > 0 {
		Offset += MaxPageSize
		goto SEARCH
	}

	return Resp, nil
}

// SearchWithKey search key and version  with key
func (d *Elasticsearch) SearchWithKey(ctx context.Context, key string) (interface{}, error) {
	ql := d.client.Search().Index(d.esConfig.DefaultIndex)

	result, err := ql.Query(elastic.NewPrefixQuery("key", key)).
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("key")).
		Size(500).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	type ksv struct {
		Key     string `json:"key,omitempty"`
		Version string `json:"version,omitempty"`
	}

	body := make([]ksv, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		ss := strings.SplitN(hit.Id, "_", 2)
		if len(ss) != 2 {
			continue
		}

		body = append(body,
			ksv{
				Key:     ss[1],
				Version: ss[0],
			})
	}

	return body, nil
}

// DeleteWithKey delete with key
func (d *Elasticsearch) DeleteWithKey(ctx context.Context, key string) error {
	ql := d.client.Search().Index(d.esConfig.DefaultIndex)

	result, err := ql.Query(elastic.NewPrefixQuery("key", key)).
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("key")).
		Size(500).
		Do(ctx)

	if err != nil {
		return err
	}

	for _, hit := range result.Hits.Hits {
		_, _ = d.client.Delete().
			Index(d.esConfig.DefaultIndex).
			Id(hit.Id).
			Do(ctx)
	}

	return nil
}

package persona

import (
	"context"
	"encoding/json"

	"git.internal.yunify.com/qxp/persona/internal/model"
	"git.internal.yunify.com/qxp/persona/internal/server/options"
	"git.internal.yunify.com/qxp/persona/pkg/config"
	"git.internal.yunify.com/qxp/persona/pkg/db"
	"git.internal.yunify.com/qxp/persona/pkg/db/elasticsearch"
	"git.internal.yunify.com/qxp/persona/pkg/misc/id2"
	"git.internal.yunify.com/qxp/persona/pkg/misc/time2"
	"git.internal.yunify.com/qxp/persona/pkg/utils"
)

// Persona inter
type Persona interface {
	UserSetValue(ctx context.Context, req *BatchSetValueReq) (*BatchSetValueResp, error)
	UserGetValue(ctx context.Context, req *BatchGetValueReq) (*BatchGetValueResp, error)
	SetValue(ctx context.Context, req *BatchSetValueReq) (*BatchSetValueResp, error)
	GetValue(ctx context.Context, req *BatchGetValueReq) (*BatchGetValueResp, error)
	CloneValue(ctx context.Context, req *CloneValueReq) (string, error)
	ExportData(ctx context.Context, req *ExportDataReq) (*ExportDataResp, error)
	ImportData(ctx context.Context, req *ImportDataReq) error

	GetDataSetByID(ctx context.Context, req *GetDataSetReq) (*GetDataSetResp, error)
	CreateDataset(ctx context.Context, req *CreateDataSetReq) (*CreateDataSetResp, error)
	UpdateDataSet(ctx context.Context, req *UpdateDataSetReq) (*UpdateDataSetResp, error)
	GetByConditionSet(c context.Context, req *GetByConditionSetReq) (*GetByConditionSetResp, error)
	DeleteDataSet(c context.Context, req *DeleteDataSetReq) (*DeleteDataSetResp, error)
}

type persona struct {
	conf    *config.Configs
	daoRepo db.BackendStorage
}

// NewPersona new
func NewPersona(conf *config.Configs, opts ...options.Options) (Persona, error) {
	dao, err := elasticsearch.NewEs(conf)
	if err != nil {
		return nil, err
	}
	return &persona{
		conf:    conf,
		daoRepo: dao,
	}, nil
}

func (p *persona) UserSetValue(ctx context.Context, req *BatchSetValueReq) (*BatchSetValueResp, error) {
	successKeys := make([]string, 0)
	failKeys := make([]string, 0)
	for _, value := range req.Keys {
		err := p.daoRepo.UserPutWithVersion(ctx, value.Version, value.Key, value.Value)
		if err != nil {
			failKeys = append(failKeys, value.Key)
		} else {
			successKeys = append(successKeys, value.Key)
		}
	}

	return &BatchSetValueResp{
		SuccessKeys: successKeys,
		FailKeys:    failKeys,
	}, nil
}

func (p *persona) UserGetValue(ctx context.Context, req *BatchGetValueReq) (*BatchGetValueResp, error) {
	result := make(map[string]string, 0)
	for _, value := range req.Keys {
		r, err := p.daoRepo.UserGetWithVersion(ctx, value.Version, value.Key)
		if err == nil && len(r) > 0 {
			result = utils.MergeMap2(result, r)
		}
	}

	return &BatchGetValueResp{
		Result: result,
	}, nil
}

func (p *persona) SetValue(ctx context.Context, req *BatchSetValueReq) (*BatchSetValueResp, error) {
	successKeys := make([]string, 0)
	failKeys := make([]string, 0)
	for _, value := range req.Keys {
		err := p.daoRepo.PutWithVersion(ctx, value.Version, value.Key, value.Value)
		if err != nil {
			failKeys = append(failKeys, value.Key)
		} else {
			successKeys = append(successKeys, value.Key)
		}
	}

	return &BatchSetValueResp{
		SuccessKeys: successKeys,
		FailKeys:    failKeys,
	}, nil
}

func (p *persona) GetValue(ctx context.Context, req *BatchGetValueReq) (*BatchGetValueResp, error) {
	result := make(map[string]string, 0)
	for _, value := range req.Keys {
		r, err := p.daoRepo.GetWithVersion(ctx, value.Version, value.Key)
		if err == nil && len(r) > 0 {
			result = utils.MergeMap2(result, r)
		}
	}

	return &BatchGetValueResp{
		Result: result,
	}, nil
}

func (p *persona) CloneValue(ctx context.Context, req *CloneValueReq) (string, error) {
	r, err := p.daoRepo.GetWithVersion(ctx, req.Key.Version, req.Key.Key)
	if err != nil {
		return "", err
	}
	if len(r) > 0 {
		var value string
		for _, v := range r {
			value = v
		}

		err := p.daoRepo.PutWithVersion(ctx, req.NewKey.Version, req.NewKey.Key, value)
		if err != nil {
			return "", err
		}

		return value, nil
	}

	return "", nil
}

func (p *persona) ExportData(ctx context.Context, req *ExportDataReq) (*ExportDataResp, error) {
	datas, err := p.daoRepo.GetWithPrefix(ctx, "app_id:"+req.AppID)
	if err != nil {
		return nil, err
	}

	return &ExportDataResp{
		AppData: datas,
	}, nil
}

func (p *persona) ImportData(ctx context.Context, req *ImportDataReq) error {
	for _, data := range req.AppData {
		err := p.daoRepo.Put(ctx, data.Key, data.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateDataset 创建数据集
func (p *persona) CreateDataset(ctx context.Context, req *CreateDataSetReq) (*CreateDataSetResp, error) {
	key := id2.GenID()
	dataset := model.DataSet{
		ID:        key,
		Name:      req.Name,
		Tag:       req.Tag,
		Type:      req.Type,
		Content:   req.Content,
		CreatedAt: time2.NowUnix(),
		DataType:  elasticsearch.TypeOfDataSet,
	}
	if err := p.daoRepo.PutData(&ctx, &key, dataset); err != nil {
		return nil, err
	}

	return &CreateDataSetResp{ID: key}, nil
}

// GetDataSetByID 根据ID获取数据集
func (p *persona) GetDataSetByID(ctx context.Context, req *GetDataSetReq) (*GetDataSetResp, error) {
	var resp GetDataSetResp
	data, _ := p.daoRepo.GetData(&ctx, &req.ID)

	//if err != nil {
	//	return nil, err
	//}
	if data == nil {
		return &resp, nil
	}
	if err := json.Unmarshal(*data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateDataSet 更新数据集
func (p *persona) UpdateDataSet(ctx context.Context, req *UpdateDataSetReq) (*UpdateDataSetResp, error) {
	if err := p.daoRepo.UpdateData(&ctx, &req.ID, req); err != nil {
		return nil, err
	}
	return &UpdateDataSetResp{}, nil
}

func (p *persona) GetByConditionSet(ctx context.Context, req *GetByConditionSetReq) (*GetByConditionSetResp, error) {
	var Resp GetByConditionSetResp
	var FilterKVs = make(map[string]interface{})
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &FilterKVs); err != nil {
		return nil, err
	}
	dataList, err := p.daoRepo.GetDataByKVs(&ctx, &FilterKVs)
	for _, d := range dataList {
		var data DataSetVo
		if err := json.Unmarshal(*d, &data); err != nil {
			return nil, err
		}
		Resp.List = append(Resp.List, &data)
	}
	return &Resp, nil
}

func (p *persona) DeleteDataSet(ctx context.Context, req *DeleteDataSetReq) (*DeleteDataSetResp, error) {
	if err := p.daoRepo.DeleteData(&ctx, &req.ID); err != nil {
		return nil, err
	}
	return &DeleteDataSetResp{}, nil
}

// ExportDataReq req
type ExportDataReq struct {
	AppID string `json:"appId" binding:"required"`
}

// ExportDataResp resp
type ExportDataResp struct {
	AppData []db.ImportReqData `json:"appData"`
}

// ImportDataReq req
type ImportDataReq struct {
	// AppID   string `json:"appId" binding:"required"`
	AppData []db.ImportReqData `json:"appData" binding:"required"`
}

// CloneValueReq req
type CloneValueReq struct {
	Key    VersionKey `json:"key" binding:"required"`
	NewKey VersionKey `json:"newKey" binding:"required"`
}

// BatchSetValueReq req
type BatchSetValueReq struct {
	Keys []VersionKeyValue `json:"keys" binding:"required"`
}

// BatchSetValueResp resp
type BatchSetValueResp struct {
	SuccessKeys []string `json:"successKeys"`
	FailKeys    []string `json:"failKeys"`
}

// BatchGetValueReq req
type BatchGetValueReq struct {
	Keys []VersionKey `json:"keys" binding:"required"`
}

// BatchGetValueResp resp
type BatchGetValueResp struct {
	Result map[string]string `json:"result"`
}

// VersionKeyValue req
type VersionKeyValue struct {
	Version string `json:"version" binding:"required"`
	Key     string `json:"key" binding:"required"`
	Value   string `json:"value" binding:"required"`
}

// VersionKey req
type VersionKey struct {
	Version string `json:"version" binding:"required"`
	Key     string `json:"key" binding:"required"`
}

// GetDataSetReq 根据ID获取数据请求
type GetDataSetReq struct {
	ID string `json:"id"`
}

// GetDataSetResp 获取数据集单条数据结构
type GetDataSetResp struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty" binding:"max=100"`
	Tag       string `json:"tag,omitempty"  binding:"max=100"`
	Type      int64  `json:"type,omitempty"`
	Content   string `json:"content,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
}

// CreateDataSetReq 新增数据集请求
type CreateDataSetReq struct {
	Name    string `json:"name" binding:"max=100"`
	Tag     string `json:"tag"  binding:"max=100"`
	Type    int64  `json:"type"`
	Content string `json:"content"`
}

// CreateDataSetResp 新增数据集响应
type CreateDataSetResp struct {
	ID string `json:"id"`
}

// UpdateDataSetReq UpdateDataSetReq
type UpdateDataSetReq struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Type    int64  `json:"type"`
	Content string `json:"content"`
}

// UpdateDataSetResp UpdateDataSetResp
type UpdateDataSetResp struct {
}

// GetByConditionSetReq 获取数据集筛选条件请求
type GetByConditionSetReq struct {
	Name  string `json:"name,omitempty"`
	Tag   string `json:"tag,omitempty"`
	Types int64  `json:"type,omitempty"`
}

// GetByConditionSetResp 获取数据集返回
type GetByConditionSetResp struct {
	List []*DataSetVo `json:"list"`
}

// DataSetVo DataSetVo
type DataSetVo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	Type      int64  `json:"type"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

// DeleteDataSetReq DeleteDataSetReq
type DeleteDataSetReq struct {
	ID string `json:"id"`
}

// DeleteDataSetResp DeleteDataSetResp
type DeleteDataSetResp struct {
}

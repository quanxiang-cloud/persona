package restful

import (
	"context"

	"git.internal.yunify.com/qxp/persona/pkg/config"
	"git.internal.yunify.com/qxp/persona/pkg/probe"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
)

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
)

// Router 路由
type Router struct {
	*probe.Probe

	c   *config.Configs
	log logr.Logger

	engine *gin.Engine

	cancel context.CancelFunc
}

// NewRouter 开启路由
func NewRouter(c *config.Configs, log logr.Logger) (*Router, error) {
	engine, err := newRouter(c, log)
	if err != nil {
		return nil, err
	}

	probe := probe.New(log)

	ctx, cancel := context.WithCancel(context.Background())
	p, err := NewPersona(ctx, c)
	if err != nil {
		return nil, err
	}
	v1 := engine.Group("/api/v1/persona")
	{
		v1.POST("/userBatchSetValue", p.userSetValue)
		v1.POST("/userBatchGetValue", p.userGetValue)

		v1.POST("/batchSetValue", p.setValue)
		v1.POST("/batchGetValue", p.getValue)

		v1.POST("/cloneValue", p.cloneValue)

		v1.POST("/app/import", p.importData)
		v1.POST("/app/export", p.exportData)
	}

	// 数据集
	smAPI := engine.Group("/api/v1/persona/dataset/m")
	{
		// 创建数据集
		smAPI.POST("/create", p.createDataSet)
		// 根据ID获取数据集
		smAPI.POST("/get", p.getDataSetByID)
		// 修改数据集
		smAPI.POST("/update", p.updateDataSet)
		// 根据条件获取结果集列表（不分页）
		smAPI.POST("/getByCondition", p.getDataSetByCondition)
		// 删除数据集
		smAPI.POST("/delete", p.deleteDataSet)

		smAPI.POST("/search/key", p.searchWithKey)
		smAPI.POST("/bulk/delete", p.deleteWithKey)
	}
	// 用户端API
	suAPI := engine.Group("/api/v1/persona/dataset/home")
	{
		// 根据ID获取数据集
		suAPI.POST("/get", p.getDataSetByIDHome)
	}
	router := &Router{
		c:      c,
		log:    log,
		Probe:  probe,
		engine: engine,
		cancel: cancel,
	}

	router.probe()
	return router, nil
}

func newRouter(c *config.Configs, log logr.Logger) (*gin.Engine, error) {
	if c.Model == "" || (c.Model != ReleaseMode && c.Model != DebugMode) {
		c.Model = ReleaseMode
	}
	gin.SetMode(c.Model)
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	return engine, nil
}

func (r *Router) probe() {
	r.engine.GET("liveness", func(c *gin.Context) {
		r.Probe.LivenessProbe(c.Writer, c.Request)
	})

	r.engine.Any("readiness", func(c *gin.Context) {
		r.Probe.ReadinessProbe(c.Writer, c.Request)
	})
}

// Run 启动服务
func (r *Router) Run() {
	r.Probe.SetRunning()
	r.engine.Run(r.c.Port)
}

// Close 关闭服务
func (r *Router) Close() {
	r.cancel()
}

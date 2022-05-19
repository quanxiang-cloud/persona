package restful

import (
	"context"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/resp"
	"git.internal.yunify.com/qxp/persona/internal/persona"
	"git.internal.yunify.com/qxp/persona/internal/server/options"
	"git.internal.yunify.com/qxp/persona/pkg/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Persona Persona
type Persona struct {
	persona persona.Persona
}

// NewPersona new persona
func NewPersona(ctx context.Context, c *config.Configs, opts ...options.Options) (*Persona, error) {
	p, err := persona.NewPersona(c)
	if err != nil {
		return nil, err
	}
	return &Persona{
		persona: p,
	}, nil
}

func (p *Persona) userSetValue(c *gin.Context) {
	req := &persona.BatchSetValueReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.persona.UserSetValue(logger.CTXTransfer(c), req)).Context(c)
}

func (p *Persona) userGetValue(c *gin.Context) {
	req := &persona.BatchGetValueReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.persona.UserGetValue(logger.CTXTransfer(c), req)).Context(c)
}

func (p *Persona) setValue(c *gin.Context) {
	req := &persona.BatchSetValueReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.persona.SetValue(logger.CTXTransfer(c), req)).Context(c)
}

func (p *Persona) getValue(c *gin.Context) {
	req := &persona.BatchGetValueReq{}

	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp.Format(p.persona.GetValue(logger.CTXTransfer(c), req)).Context(c)
}

func (p *Persona) cloneValue(c *gin.Context) {
	req := &persona.CloneValueReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.persona.CloneValue(logger.CTXTransfer(c), req)).Context(c)
}

func (p *Persona) exportData(c *gin.Context) {
	req := &persona.ExportDataReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.persona.ExportData(logger.CTXTransfer(c), req)).Context(c)
}

func (p *Persona) importData(c *gin.Context) {
	req := &persona.ImportDataReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(nil, p.persona.ImportData(logger.CTXTransfer(c), req)).Context(c)
}

// createDataSet 创建数据集(管理端)
func (p *Persona) createDataSet(c *gin.Context) {
	req := &persona.CreateDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.persona.CreateDataset(logger.CTXTransfer(c), req)).Context(c)
}

// getDataSetByID 根据ID得到数据集(管理端)
func (p *Persona) getDataSetByID(c *gin.Context) {
	req := &persona.GetDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.persona.GetDataSetByID(logger.CTXTransfer(c), req)).Context(c)
}

// UpdateDataSet 修改数据集(管理端)
func (p *Persona) updateDataSet(c *gin.Context) {
	req := &persona.UpdateDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.persona.UpdateDataSet(logger.CTXTransfer(c), req)).Context(c)

}

// getDataSetByIDHome 根据ID得到数据集（用户端）
func (p *Persona) getDataSetByIDHome(c *gin.Context) {
	req := &persona.GetDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.persona.GetDataSetByID(logger.CTXTransfer(c), req)).Context(c)
}

// getDataSetByCondition 根据条件获取数据集(管理端)
func (p *Persona) getDataSetByCondition(c *gin.Context) {
	req := &persona.GetByConditionSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.persona.GetByConditionSet(logger.CTXTransfer(c), req)).Context(c)
}

// deleteDataSet 删除数据集
func (p *Persona) deleteDataSet(c *gin.Context) {
	req := &persona.DeleteDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.persona.DeleteDataSet(logger.CTXTransfer(c), req)).Context(c)
}

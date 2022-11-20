package project

import (
	"icarus/utils"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type Controller struct {
	Ctx     iris.Context
	Service ProjectService
}

func (c *Controller) PostCreate() mvc.Result {
	var project map[string]interface{}
	if err := c.Ctx.ReadJSON(&project); err != nil {
		return utils.RestfulResponse(4000, err.Error(), map[string]string{})
	}
	result, err := c.Service.Create(project)
	if err != nil {
		return utils.RestfulResponse(4022, err.Error(), map[string]string{})
	}
	return utils.RestfulResponse(2000, "create success", result)
}

func (c *Controller) PutUpdate() mvc.Result {
	params := make(map[string]interface{})
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return utils.RestfulResponse(4000, err.Error(), map[string]string{})
	}
	updateProject, err := c.Service.Update(params)
	if err != nil {
		return utils.RestfulResponse(4000, err.Error(), map[string]string{})
	}
	return utils.RestfulResponse(2000, "update success", updateProject)
}

func (c *Controller) GetBy(id uint16) mvc.Result {
	selectProject, err := c.Service.GetbyID(id)
	if err != nil {
		return utils.RestfulResponse(4004, err.Error(), map[string]string{})
	}
	return utils.RestfulResponse(2000, "found", selectProject)
}

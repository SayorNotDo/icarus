package project

import (
	"icarus/user"
	"icarus/utils"
	"log"

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
	result, code, err := c.Service.Create(project)
	if err != nil {
		return utils.RestfulResponse(code, err.Error(), map[string]string{})
	}
	return utils.RestfulResponse(code, "Create Success", result)
}

func (c *Controller) PutUpdate() mvc.Result {
	params := make(map[string]interface{})
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return utils.RestfulResponse(4000, err.Error(), map[string]string{})
	}
	updateProject, code, err := c.Service.Update(params)
	if err != nil {
		return utils.RestfulResponse(code, err.Error(), map[string]string{})
	}
	return utils.RestfulResponse(code, "Update Success", updateProject)
}

func (c *Controller) GetBy(id uint16) mvc.Result {
	selectProject, code, err := c.Service.GetbyID(id)
	if err != nil {
		return utils.RestfulResponse(code, err.Error(), map[string]string{})
	}
	return utils.RestfulResponse(code, "Project Found", selectProject)
}

// TODO: Get all projects interface complete
func (c *Controller) Get() mvc.Result {
	uid, username := user.ParseUserinfo(c.Ctx)
	// TODO: check if username have authority to get all projects
	log.Printf("uid: %v, username: %s", uid, username)
	return utils.RestfulResponse(2000, "", map[string]string{})
}

// TODO: Delete interface complete
func (c *Controller) DeleteBy(pid uint16) mvc.Result {
	return utils.RestfulResponse(2000, "nil", map[string]string{})
}

// TODO:ProjectMember Insert interface complete
func (c *Controller) PostMember() mvc.Result {
	return utils.RestfulResponse(2000, "", map[string]string{})
}

// TODO:ProjectMember Delete interface complete
func (c *Controller) DeleteMemberBy(id int16) mvc.Result {
	return utils.RestfulResponse(2000, "", map[string]string{})
}

// TODO: ProjectMember Update interface complete
func (c *Controller) PutMember() mvc.Result {
	return utils.RestfulResponse(2000, "", map[string]string{})
}

// TODO: ProjectMember Get interface complete
func (c *Controller) GetMember() mvc.Result {
	return utils.RestfulResponse(2000, "", map[string]string{})
}

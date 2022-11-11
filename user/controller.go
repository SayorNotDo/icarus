package user

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"icarus/utils"
	"log"
)

type Controller struct {
	Ctx     iris.Context
	Service UserService
	Session *sessions.Session
}

const userIDKey = "UID"

func (c *Controller) getCurrentUserID() int64 {
	userID := c.Session.GetInt64Default(userIDKey, 0)
	return userID
}

func (c *Controller) isLoggedIn() bool {
	return c.getCurrentUserID() > 0
}

func (c *Controller) logout() {
	c.Session.Destroy()
}

//// Get v1/api/user
//func (c *Controller) Get() (results []User) {
//	return c.Service.GetAll()
//}

//// GetBy v1/api/user/{uid:int64}
//func (c *Controller) GetBy(uid int64) (u User, found bool) {
//	return c.Service.GetByID(uid) // throw 404 if not found
//}

// PostRegister v1/api/user/register
func (c *Controller) PostRegister() mvc.Result {
	// get post data form json
	// username, password, email, phone
	var params map[string]interface{}
	// return while read json occur error
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return utils.RestfulResponse(5000, err.Error(), map[string]string{})
	}
	username := params["username"]
	password := params["password"]
	email := params["email"]
	phone := params["phone"]
	u, err := c.Service.Create(password.(string), User{
		Username: username.(string),
		Email:    email.(string),
		Phone:    phone.(string),
	})
	if err != nil {
		return utils.RestfulResponse(5000, err.Error(), map[string]string{})
	}
	c.Session.Set(userIDKey, u.UID)
	return utils.RestfulResponse(2000, "user register success!", u)
}

// GetRegister v1/api/user/register
func (c *Controller) GetRegister() mvc.Result {
	if c.isLoggedIn() {
		c.logout()
	}
	return mvc.Response{
		Text: "Register Page",
	}
}

// PostLogin v1/api/user/login
func (c *Controller) PostLogin() mvc.Result {
	var params map[string]interface{}
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return utils.RestfulResponse(2000, err.Error(), map[string]string{})
	}
	username := params["username"]
	password := params["password"]
	Token, err := c.Service.Login(username.(string), password.(string))
	if err != nil {
		return utils.RestfulResponse(2001, err.Error(), map[string]string{})
	}
	return utils.RestfulResponse(2000, "login success!", map[string]string{"token": string(Token)})
}

// PostLogout v1/api/user/logout
func (c *Controller) PostLogout() mvc.Result {
	return utils.RestfulResponse(2000, "user has logout!", map[string]string{})
}

// PutUpdate v1/api/user/update
func (c *Controller) PutUpdate() mvc.Result {
	var params map[string]interface{}
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return utils.RestfulResponse(2000, err.Error(), map[string]string{})
	}
	log.Printf("json: %v", params)
	return mvc.Response{
		Text: "Update success",
	}
}

func (u User) IsValid() bool {
	return u.UID > 0
}
func (u User) Dispatch(ctx context.Context) {
	if !u.IsValid() {
		ctx.NotFound()
		return
	}
	ctx.JSON(u, context.JSON{Indent: " "})
}

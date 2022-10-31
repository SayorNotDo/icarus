package user

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
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

// Get v1/api/user
func (c *Controller) Get() (results []User) {
	return c.Service.GetAll()
}

// GetBy v1/api/user/{uid:int64}
func (c *Controller) GetBy(uid int64) (u User, found bool) {
	return c.Service.GetByID(uid) // throw 404 if not found
}

// PostRegister v1/api/user/register
func (c *Controller) PostRegister() mvc.Result {
	// get post data form json
	// username, password, email, phone
	var params map[string]interface{}
	// return while read json occur error
	if err := c.Ctx.ReadJSON(&params); err != nil {
		c.Ctx.JSON(iris.Map{
			"code":    5000,
			"message": err.Error(),
			"data":    "",
		})
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
		c.Ctx.JSON(iris.Map{
			"code":    2000,
			"message": err.Error(),
			"data":    "",
		})
	}
	c.Session.Set(userIDKey, u.UID)
	c.Ctx.JSON(iris.Map{
		"code":    2000,
		"message": "User register success!",
		"data":    "{}",
	})
	return nil
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
		return mvc.Response{
			Text: err.Error(),
		}
	}
	//username := params["username"]
	//password := params["password"]
	return mvc.Response{
		Text: fmt.Sprintf("%v, %v", params["username"], params["password"]),
	}
}

// PostAuthenticate v1/api/user/authenticate
func (c *Controller) PostAuthenticate() mvc.Result {
	var params map[string]interface{}
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return mvc.Response{
			Text: err.Error(),
		}
	}
	return mvc.Response{
		Text: "--------------authenticate interface--------------",
	}
}

// PostLogout http://localhost:8080/v1/api/user/logout
func (c *Controller) PostLogout() mvc.Result {
	return mvc.Response{
		Text: "Logout success!",
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

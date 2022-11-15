package user

import (
	"icarus/utils"
	"log"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
)

type Controller struct {
	Ctx     iris.Context
	Service UserService
}

// func (c *Controller) getCurrentUserID() int64 {
// 	userID := c.Session.GetInt64Default(userIDKey, 0)
// 	return userID
// }

// func (c *Controller) isLoggedIn() bool {
// 	return c.getCurrentUserID() > 0
// }

// func (c *Controller) logout() {
// 	c.Session.Destroy()
// }

// PostRegister v1/api/user/register
func (c *Controller) PostRegister() mvc.Result {
	// get post data form json
	// username, password, email, phone
	var params map[string]string
	// return while read json occur error
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return utils.RestfulResponse(5000, err.Error(), map[string]string{})
	}
	u, err := c.Service.Create(params)
	if err != nil {
		return utils.RestfulResponse(5000, err.Error(), map[string]string{})
	}
	return utils.RestfulResponse(2000, "user register success!", u)
}

// GetRegister v1/api/user/register
// func (c *Controller) GetRegister() mvc.Result {
// 	if c.isLoggedIn() {
// 		c.logout()
// 	}
// 	// redirect to register page
// 	return mvc.Response{
// 		Text: "Register Page",
// 	}
// }

// PostLogin v1/api/user/login
func (c *Controller) PostLogin() mvc.Result {
	var params map[string]interface{}
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return utils.RestfulResponse(2000, err.Error(), map[string]string{})
	}
	username := params["username"]
	password := params["password"]
	Token, refreshToken, err := c.Service.Login(username.(string), password.(string))
	if err != nil {
		return utils.RestfulResponse(2001, err.Error(), map[string]string{})
	}
	return utils.RestfulResponse(2000, "login success!", map[string]string{"token": Token, "refreshToken": refreshToken})
}

// PostLogout v1/api/user/logout
func (c *Controller) PostLogout() mvc.Result {
	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	res := token.Claims.(jwt.MapClaims)
	uid := res["uid"]
	username := res["username"]
	log.Printf("uid: %v, username: %v", uid, username)

	return utils.RestfulResponse(2000, "user has logout!", map[string]string{})
}

// PutUpdate v1/api/user/update
func (c *Controller) PutUpdate() mvc.Result {
	Authorization := c.Ctx.GetHeader("Authorization")
	log.Printf("get Authorization parameter: %s", Authorization)
	var params map[string]string
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return utils.RestfulResponse(2000, err.Error(), map[string]string{})
	}
	log.Printf("json: %v", params)
	if _, err := c.Service.Update(params); err != nil {
		return mvc.Response{
			Text: "update failed",
		}
	}
	return mvc.Response{
		Text: "update success",
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

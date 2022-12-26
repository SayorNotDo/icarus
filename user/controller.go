package user

import (
	. "icarus/utils"
	"log"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
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

// GetUserBy v1/api/user/<string:username>
func (c *Controller) GetUserBy(username string) (User, error) {
	panic("not implemented")
}

// PostRegister v1/api/user/register
func (c *Controller) PostRegister() mvc.Result {
	var params map[string]string
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return Response(5000, err.Error(), map[string]string{})
	}
	if params["username"] == "" || params["password"] == "" {
		return Response(iris.StatusBadRequest, "username or password can not be empty", map[string]string{})
	} else if params["email"] == "" {
		return Response(iris.StatusBadRequest, "email can not be empty", map[string]string{})
	} else if params["phone"] == "" {
		return Response(iris.StatusBadRequest, "phone number can not be empty", map[string]string{})
	}
	hashed, err := GeneratePassword(params["password"])
	if err != nil {
		return Response(iris.StatusInternalServerError, "password hashed error", map[string]string{})
	}
	createInfo := map[string]string{
		"username": params["username"],
		"email":    params["email"],
		"phone":    params["phone"],
	}
	if status := c.Service.Create(createInfo, hashed); status.err != nil {
		return Response(status.statusCode, status.err.Error(), map[string]string{})
	}
	return Response(iris.StatusCreated, "register success!", map[string]string{})
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
		return Response(200, err.Error(), map[string]string{})
	}
	username := params["username"]
	password := params["password"]
	Token, refreshToken, err := c.Service.Login(username.(string), password.(string))
	if err != nil {
		return Response(200, err.Error(), map[string]string{})
	}
	return Response(200, "login success!",
		iris.Map{
			"accessToken":  Token,
			"refreshToken": refreshToken,
			"tokenType":    "Bearer",
		})
}

// PostLogout v1/api/user/logout
func (c *Controller) PostLogout() mvc.Result {
	uid, username := ParseUserinfo(c.Ctx)
	log.Printf("uid: %v, username: %v", uid, username)
	params := map[string]interface{}{
		"username":      username,
		"refresh_token": ""}
	c.Service.Logout(params)
	return Response(2000, "user has logout!", map[string]string{})
}

func (c *Controller) PostAuthenticate() mvc.Result {
	var params map[string]string
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return Response(2000, err.Error(), map[string]string{})
	}
	token, err := jwt.FromAuthHeader(c.Ctx)
	if err != nil {
		return Response(2000, err.Error(), map[string]string{})
	}
	newToken, newRefreshToken, err := c.Service.Authenticate(token, params["refreshToken"])
	if err != nil {
		return Response(5000, err.Error(), map[string]string{})
	}
	return Response(2000, "authenticate success", iris.Map{"accessToken": newToken, "refreshToken": newRefreshToken, "tokenType": "Bearer"})
}

// PutUpdate v1/api/user/update
func (c *Controller) PutUpdate() mvc.Result {
	uid, username := ParseUserinfo(c.Ctx)
	log.Printf("%v, %v", uid, username)
	var params map[string]interface{}
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return Response(5000, err.Error(), map[string]string{})
	}
	log.Printf("json: %v", params)
	if _, err := c.Service.Update(User{Username: username, UID: uid}, params); err != nil {
		return Response(5000, err.Error(), map[string]string{})
	}
	return Response(2000, "update success", map[string]string{})
}

func (c *Controller) DeleteBy(id uint32) mvc.Result {
	uid, username := ParseUserinfo(c.Ctx)
	log.Printf("%v, %v", uid, username)
	isDelete := c.Service.DeleteByID(id)
	log.Println(isDelete)
	return Response(2000, "delete success", map[string]string{})
}

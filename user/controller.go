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

// Get v1/api/user
func (c *Controller) Get() mvc.Result {
	user := c.Service.GetUserInfo(c.Ctx)
	if user == nil {
		return ResponseError(iris.StatusInternalServerError, "error get user info")
	}
	return ResponseSuccess("success!", user)
}

func (c *Controller) PostRegister() mvc.Result {
	var params map[string]string
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return ResponseError(iris.StatusInternalServerError, err.Error())
	}
	if params["username"] == "" || params["password"] == "" {
		return ResponseError(iris.StatusBadRequest, "username or password can not be empty")
	} else if params["email"] == "" {
		return ResponseError(iris.StatusBadRequest, "email can not be empty")
	} else if params["phone"] == "" {
		return ResponseError(iris.StatusBadRequest, "phone number can not be empty")
	}
	hashed, err := generatePassword(params["password"])
	if err != nil {
		return ResponseError(iris.StatusInternalServerError, "password hashed error")
	}
	createInfo := map[string]string{
		"username": params["username"],
		"email":    params["email"],
		"phone":    params["phone"],
	}
	if status := c.Service.Create(createInfo, hashed); status.err != nil {
		return ResponseError(status.statusCode, status.err.Error())
	}
	return Response(iris.StatusCreated, "register success!", map[string]string{})
}

func (c *Controller) PostLogin() mvc.Result {
	params := make(map[string]string)
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return ResponseError(iris.StatusInternalServerError, err.Error())
	}
	if params["username"] == "" || params["password"] == "" {
		return ResponseError(iris.StatusBadRequest, "username or password can not be empty")
	}
	tokenMap, status := c.Service.Login(params["username"], params["password"])
	if status.err != nil {
		return ResponseError(status.statusCode, status.err.Error())
	}
	return Response(iris.StatusOK, "login success!",
		tokenMap)
}

func (c *Controller) PostLogout() mvc.Result {
	uid, username := ParseUserinfo(c.Ctx)
	log.Printf("uid: %v, username: %v", uid, username)
	params := map[string]interface{}{
		"username":      username,
		"refresh_token": "",
	}
	if err := c.Service.Logout(params); err != nil {
		return Response(iris.StatusInternalServerError, "error", map[string]string{})
	}
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
	newToken, newRefreshToken, status := c.Service.Authenticate(token, params["refreshToken"])
	if status.err != nil {
		return Response(status.statusCode, status.err.Error(), map[string]string{})
	}
	return Response(iris.StatusOK, "authenticate success", iris.Map{"accessToken": newToken, "refreshToken": newRefreshToken, "tokenType": "Bearer"})
}

func (c *Controller) PostAuthorize() mvc.Result {
	params := make(map[string]string)
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return Response(iris.StatusInternalServerError, err.Error(), map[string]string{})
	}
	if params["username"] == "" || params["password"] == "" {
		return Response(iris.StatusBadRequest, "username or password can not be empty", map[string]string{})
	}
	authorizeToken, status := c.Service.Authorize(params["username"], params["password"])
	if status.err != nil {
		return Response(status.statusCode, status.err.Error(), map[string]string{})
	}
	return Response(iris.StatusOK, "authorize success", iris.Map{
		"authorizeToken": authorizeToken,
		"tokenType":      "Bearer",
	})
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
	if _, err := c.Service.Update(&User{Username: username, UID: uid}, params); err != nil {
		return Response(5000, err.Error(), map[string]string{})
	}
	return Response(2000, "update success", map[string]string{})
}

// DeleteBy TODO: implement interface
func (c *Controller) DeleteBy(id uint32) mvc.Result {
	uid, username := ParseUserinfo(c.Ctx)
	log.Printf("%v, %v", uid, username)
	//if isAdmin := validateAdministrator(uid); !isAdmin {
	//	return Response(iris.StatusForbidden, "you don't have permission", map[string]string{})
	//}
	c.Service.DeleteByID(id)
	return Response(2000, "delete success", map[string]string{})
}

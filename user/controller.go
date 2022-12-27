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
	uid, username := ParseUserinfo(c.Ctx)
	user, found := c.Service.GetUserInfo(uid, username)
	if !found {
		return Response(iris.StatusInternalServerError, "error get user info", map[string]string{})
	}
	return Response(iris.StatusOK, "success!", user)
}

func (c *Controller) PostRegister() mvc.Result {
	var params map[string]string
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return Response(500, err.Error(), map[string]string{})
	}
	if params["username"] == "" || params["password"] == "" {
		return Response(iris.StatusBadRequest, "username or password can not be empty", map[string]string{})
	} else if params["email"] == "" {
		return Response(iris.StatusBadRequest, "email can not be empty", map[string]string{})
	} else if params["phone"] == "" {
		return Response(iris.StatusBadRequest, "phone number can not be empty", map[string]string{})
	}
	hashed, err := generatePassword(params["password"])
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

func (c *Controller) PostLogin() mvc.Result {
	params := make(map[string]string)
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return Response(iris.StatusInternalServerError, err.Error(), map[string]string{})
	}
	if params["username"] == "" || params["password"] == "" {
		return Response(iris.StatusBadRequest, "username or password can not be empty", map[string]string{})
	}
	accessToken, refreshToken, status := c.Service.Login(params["username"], params["password"])
	if status.err != nil {
		return Response(status.statusCode, status.err.Error(), map[string]string{})
	}
	return Response(iris.StatusOK, "login success!",
		iris.Map{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
			"tokenType":    "Bearer",
		})
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
	if _, err := c.Service.Update(User{Username: username, UID: uid}, params); err != nil {
		return Response(5000, err.Error(), map[string]string{})
	}
	return Response(2000, "update success", map[string]string{})
}

// DeleteBy TODO: implement interface
func (c *Controller) DeleteBy(id uint32) mvc.Result {
	uid, username := ParseUserinfo(c.Ctx)
	log.Printf("%v, %v", uid, username)
	if isAdmin := validateAdministrator(uid); !isAdmin {
		return Response(iris.StatusForbidden, "you don't have permission", map[string]string{})
	}
	isDelete := c.Service.DeleteByID(id)
	log.Println(isDelete)
	return Response(2000, "delete success", map[string]string{})
}

package user

import (
	"icarus/static"
	"strings"
	"time"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

var (
	secret = []byte("signature_hmac_secret_shared_key")
	// encryption = []byte("signature_hmac_secret_shared_key")
)

var Signer = jwt.New(jwt.Config{
	ErrorHandler: func(ctx iris.Context, err error) {
		if err == nil {
			return
		}
		ctx.JSON(iris.Map{
			"code":    401,
			"message": err.Error(),
			"data":    map[string]string{},
		})
	},
	// get jwt from Authorization in Request Header
	Extractor: jwt.FromAuthHeader,

	Expiration: true,

	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	},

	SigningMethod: jwt.SigningMethodHS256,
})

func generateAccessToken(username string, uid uint32) (token string, err error) {
	now := time.Now()
	generateToken := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"uid":      uid,
		"iat":      now.Unix(),
		"exp":      now.Add(150 * time.Hour).Unix(),
	})

	tokenString, err := generateToken.SignedString(secret)
	return tokenString, err
}

func generateRefreshToken(token string) (refreshToken string, err error) {
	generateToken := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"accessToken": token,
	})
	refreshTokenString, err := generateToken.SignedString(secret)
	return refreshTokenString, err
}

func authorizeToken(username string, uid uint32) (token string, err error) {
	now := time.Now()
	generateToken := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"x-mode":   true,
		"uid":      uid,
		"iat":      now.Unix(),
		"exp":      now.Add(30 * 24 * time.Hour).Unix(),
	})
	tokenString, err := generateToken.SignedString(secret)
	return tokenString, err
}

func AuthenticatedHandler(ctx iris.Context) {
	if res := exceptAuthenticate(ctx.Path()); res != false {
		if err := Signer.CheckJWT(ctx); err != nil {
			Signer.Config.ErrorHandler(ctx, err)
			return
		}
	}
	ctx.Next()
}

func exceptAuthenticate(ctxPath string) bool {
	relativePath := strings.TrimPrefix(ctxPath, static.RoutePrefix)
	switch relativePath {
	case "/user/register":
		return false
	case "/user/login":
		return false
	case "/user/authorize":
		return false
	default:
		return true
	}
}

func ParseUserinfo(ctx iris.Context) (uid uint32, username string) {
	token := ctx.Values().Get("jwt").(*jwt.Token)
	res := token.Claims.(jwt.MapClaims)
	uid = uint32(res["uid"].(float64))
	username = res["username"].(string)
	return
}

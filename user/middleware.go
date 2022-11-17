package user

import (
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
			"code":    4001,
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

func generateAccessToken(username string, uid int64) (token string, err error) {
	now := time.Now()
	generateToken := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"uid":      uid,
		"iat":      now.Unix(),
		"exp":      now.Add(15 * time.Hour).Unix(),
	})

	tokenString, err := generateToken.SignedString(secret)
	return tokenString, err
}

func generateRefreshToken(token string) (refreshToken string, err error) {
	generateToken := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": token,
	})
	refreshTokenString, err := generateToken.SignedString(secret)
	return refreshTokenString, err
}

func AuthenticatedHandler(ctx iris.Context) {
	if err := Signer.CheckJWT(ctx); err != nil {
		Signer.Config.ErrorHandler(ctx, err)
		return
	}
	ctx.Next()
}

func parseUserinfo(ctx iris.Context) (uid int64, username string) {
	token := ctx.Values().Get("jwt").(*jwt.Token)
	res := token.Claims.(jwt.MapClaims)
	uid = int64(res["uid"].(float64))
	username = res["username"].(string)
	return
}

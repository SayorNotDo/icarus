package user

import (
	"icarus/utils"
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
		ctx.StopExecution()
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(utils.RestfulResponse(501, err.Error(), map[string]string{}))
	},
	// get jwt from Authorization in Request Header
	Extractor: jwt.FromAuthHeader,

	Expiration: true,

	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	},

	SigningMethod: jwt.SigningMethodHS256,
})

func generateToken(username string, uid int64) (token string, err error) {
	now := time.Now()
	generateToken := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"uid":      uid,
		"iat":      now.Unix(),
		"exp":      now.Add(15 * time.Minute).Unix(),
	})

	tokenString, err := generateToken.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthenticatedHandler(ctx iris.Context) {
	if err := Signer.CheckJWT(ctx); err != nil {
		Signer.Config.ErrorHandler(ctx, err)
		return
	}
	ctx.Next()
}

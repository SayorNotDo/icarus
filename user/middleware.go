package user

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"time"
)

var (
	secret     = []byte("signature_hmac_secret_shared_key")
	encryption = []byte("signature_hmac_secret_shared_key")
)

type UserClaims struct {
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
}

var Signer *jwt.Signer

func init() {
	Signer = jwt.NewSigner(jwt.HS256, secret, 10*time.Minute)

	// Enable payload encryption
	Signer.WithEncryption(encryption, nil)

	verifier := jwt.NewVerifier(jwt.HS256, secret)

	verifier.WithDefaultBlocklist()
	//verifyMiddleware := verifier.Verify(func() interface{} {
	//	return new(fooClaims)
	//})
}

func generateToken(signer *jwt.Signer, username string, uid int64) (token []byte, err error) {
	claims := UserClaims{Username: username, Uid: uid}
	token, err = signer.Sign(claims)
	if err != nil {
		return nil, err
	}
	return
}

func protected(ctx iris.Context) {
	claims := jwt.Get(ctx).(*UserClaims)
	standardClaims := jwt.GetVerifiedToken(ctx).StandardClaims
	expiresAtString := standardClaims.ExpiresAt().Format(ctx.Application().ConfigurationReadOnly().GetTimeFormat())
	timeLeft := standardClaims.Timeleft()

	ctx.Writef("foo=%s\nexpires at: %s\ntime left:%s\n", claims.Username, expiresAtString, timeLeft)
}

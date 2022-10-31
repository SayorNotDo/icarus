package api

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/kataras/iris/v12/middleware/requestid"
	"icaru/database/mariadb/sql"
	"time"
)

func Router(db sql.Database, secret string) func(party iris.Party) {
	return func(r iris.Party) {
		r.Use(requestid.New())

		signer := jwt.NewSigner(jwt.HS256, secret, 15*time.Minute)
		r.Get("/token", writeToken(signer))

		verify := jwt.NewVerifier(jwt.HS256, secret).Verify(nil)
		r.Use(verify)
	}
}

func writeToken(signer *jwt.Signer) iris.Handler {
	return func(ctx iris.Context) {
		claims := jwt.Claims{
			Issuer:   "https://iris-go.com",
			Audience: []string{requestid.Get(ctx)},
		}

		token, err := signer.Sign(claims)
		if err != nil {
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}
		ctx.Write(token)
	}
}

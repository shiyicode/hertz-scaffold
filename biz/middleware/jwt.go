package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/jwt"
	"github.com/three-body/hertz-scaffold/biz/hmodel/user"
	"github.com/three-body/hertz-scaffold/biz/logic"
)

var (
	JwtMiddleware *jwt.HertzJWTMiddleware
	IdentityKey   = "identityKey"
)

func init() {
	var err error
	JwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:         "hertz jwt",
		Key:           []byte("secret key"),
		Timeout:       time.Hour * 24 * 7,
		MaxRefresh:    time.Hour * 24 * 7,
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		IdentityKey:   IdentityKey,
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return claims
		},
		// triggered when login succeed. save token data
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*user.User); ok {
				return jwt.MapClaims{
					"uid":      v.UID,
					"nickname": v.Nickname,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator:  authenticator,
		LoginResponse:  loginRespHandler,
		LogoutResponse: logoutRespHandler,
		Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
			if v, ok := data.(jwt.MapClaims); ok && v["uid"] != "" {
				return true
			}
			return false
		},
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			hlog.CtxErrorf(ctx, "jwt biz err = %+v", e.Error())
			return e.Error()
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, utils.H{
				"code":    code,
				"message": message,
			})
		},
	})
	if err != nil {
		panic(fmt.Errorf("init JWT failed: %w", err))
	}
}

func authenticator(ctx context.Context, c *app.RequestContext) (interface{}, error) {
	var req user.LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		return nil, errors.NewPublic(err.Error()).SetMeta(jwt.ErrFailedAuthentication)
	}
	if err := req.IsValid(); err != nil {
		return nil, errors.NewPublic(err.Error()).SetMeta(jwt.ErrMissingLoginValues)
	}

	user, err := logic.NewUserLogic(ctx, c).Login(&req)
	if err != nil {
		return nil, errors.NewPublic(err.Error()).SetMeta(jwt.ErrFailedAuthentication)
	}
	return user, nil
}

// loginRespHandler .
func loginRespHandler(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
	if code == http.StatusOK {
		c.JSON(http.StatusOK, user.LoginResponse{
			Token:  token,
			Expire: expire.Format(time.RFC3339),
		})
	} else {
		c.JSON(code, user.LogoutResponse{})
	}
}

// logoutRespHandler .
func logoutRespHandler(ctx context.Context, c *app.RequestContext, code int) {
	if code == http.StatusOK {
		c.JSON(http.StatusOK, user.LogoutResponse{})
	} else {
		c.JSON(code, user.LogoutResponse{})
	}
}

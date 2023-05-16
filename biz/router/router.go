package router

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/pprof"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
	"github.com/three-body/hertz-scaffold/biz/handler"
	"github.com/three-body/hertz-scaffold/biz/middleware"
	"github.com/three-body/hertz-scaffold/config"
)

func Register(h *server.Hertz) {
	h.Use(middleware.RootMw()...)

	GeneratedRegister(h)
	customizedRegister(h)
}

func customizedRegister(h *server.Hertz) {
	if config.GetConf().Hertz.EnablePprof {
		pprof.Register(h)
	}

	if config.GetConf().Swagger.Enable {
		h.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler,
			swagger.URL("/swagger/doc.json"),
			swagger.PersistAuthorization(true),
		))
		h.NoRoute(func(ctx context.Context, c *app.RequestContext) {
			if strings.HasPrefix(string(c.Path()), "/swagger/") {
				c.Next(ctx)
			} else {
				c.Redirect(consts.StatusFound, []byte("/swagger/index.html"))
			}
		})
	}

	h.GET("/ping", handler.Ping)

	// your code ...

}

package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/common/bytebufferpool"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/logger/accesslog"
	"github.com/hertz-contrib/requestid"
	"github.com/three-body/hertz-scaffold/config"
)

func RootMw() []app.HandlerFunc {
	mws := []app.HandlerFunc{
		recovery.Recovery(),
		requestid.New(),
	}
	if config.GetConf().Hertz.EnableAccessLog {
		accesslog.Tags["reqID"] = func(ctx context.Context, c *app.RequestContext, buf *bytebufferpool.ByteBuffer) (int, error) {
			return buf.WriteString(requestid.Get(c))
		}
		mws = append(mws, accesslog.New(
			accesslog.WithTimeFormat(time.RFC3339Nano),
			accesslog.WithTimeInterval(time.Millisecond),
			accesslog.WithAccessLogFunc(hlog.CtxInfof),
			// accesslog.WithFormat("[${reqID} - ${time}] ${status} ${latency} - ${method} ${path} ${queryParams} - ${ip} ${referer}"),
			accesslog.WithFormat("status=${status} latency=${latency} - ${method} ${path} ${queryParams} - ${ip} ${referer}"),
		))
	}

	return mws
}

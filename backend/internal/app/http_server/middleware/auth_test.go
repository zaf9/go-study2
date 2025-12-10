package middleware

import (
	"fmt"
	"testing"
	"time"

	appjwt "go-study2/internal/pkg/jwt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func TestAuthMiddleware(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := appjwt.Configure(appjwt.Options{
			Secret:            "0123456789abcdef",
			AccessTokenExpiry: time.Hour,
		})
		t.AssertNil(err)

		token, err := appjwt.GenerateAccessToken(42)
		t.AssertNil(err)

		s := g.Server(guid.S())
		addr := "127.0.0.1:18090"
		s.SetAddr(addr)
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.Middleware(Auth)
			group.GET("/ping", func(r *ghttp.Request) {
				uid := r.GetCtxVar("user_id").Int64()
				r.Response.WriteJson(g.Map{"uid": uid})
			})
		})

		go s.Start()
		defer s.Shutdown()

		time.Sleep(50 * time.Millisecond)

		ctx := gctx.New()
		client := g.Client()
		client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
		resp, err := client.Get(ctx, fmt.Sprintf("http://%s/ping", addr))
		t.AssertNil(err)
		defer resp.Close()

		result := resp.ReadAllString()
		t.AssertIN(`"uid":42`, result)
	})
}

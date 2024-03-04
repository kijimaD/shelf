package routers

import (
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

func NewRouter() *gin.Engine {
	gin.DefaultWriter = io.Discard

	r := gin.Default()
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: GetSlogLogLevel(Config.LogLevel),
	}))
	l = l.With("logtype", "resplog")
	slog.SetDefault(l)
	config := sloggin.Config{
		DefaultLevel:       slog.LevelInfo,
		ClientErrorLevel:   slog.LevelWarn,
		ServerErrorLevel:   slog.LevelError,
		WithUserAgent:      true,
		WithRequestID:      true,
		WithRequestBody:    true,
		WithRequestHeader:  false,
		WithResponseBody:   true,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,
	}
	r.Use(sloggin.NewWithConfig(l, config))

	r.LoadHTMLGlob("./templates/*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{},
		)
	})
	r.Static("/file", ".")

	return r
}

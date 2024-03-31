package routers

import (
	"log"

	"github.com/caarlos0/env/v10"
)

var Config AppConfig

type AppConfig struct {
	// サーバを起動するホスト.
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	// サーバを起動するポート.
	Port int `env:"PORT" envDefault:"8020"`
	// サーバを起動するアドレス.
	Address string `env:"ADDRESS,expand" envDefault:"$HOST:${PORT}"`

	// テストなどで余計な標準出力がされないようにする.ログレベルを"error"にしても、異常系のエラーなどは表示されてしまうので、それらを一切表示しないために使う
	SilentMode bool `env:"SQUALL_SILENT_MODE" envDefault:"false"`
	// ログレベル。この設定レベル以下のログは表示・保存しない。
	// 多 <- "info" | "warn" | "error" -> 少
	LogLevel LogLevel `env:"SQUALL_LOG_LEVEL" envDefault:"info"`
	// ファイルサーブのベースとなるディレクトリ
	ServeBase string `env:"SERVE_BASE" envDefault:"."`
}

func init() {
	if err := env.Parse(&Config); err != nil {
		log.Fatal(err)
	}
}

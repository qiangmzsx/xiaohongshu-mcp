package main

import (
	"context"
	"flag"

	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/browser"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
	"github.com/xpzouying/xiaohongshu-mcp/cookies"
)

// version 构建版本号，发布时通过 -ldflags "-X main.version=vX.Y.Z" 注入。
var version = "dev"

func main() {
	var (
		headless  bool
		binPath   string // 浏览器二进制文件路径
		port      string
		transport string // 传输协议：http 或 stdio
	)
	flag.BoolVar(&headless, "headless", true, "是否无头模式")
	flag.StringVar(&port, "port", ":18060", "端口")
	flag.StringVar(&transport, "transport", "http", "传输协议: http 或 stdio")
	flag.Parse()

	logrus.Infof("xiaohongshu-mcp version: %s", version)

	// 只用内置浏览器。启动时就备好，缺它直接退出，不拖到第一个请求才失败。
	binPath, err := browser.EnsureBrowser()
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	logrus.Infof("using browser binary: %s", binPath)

	configs.InitHeadless(headless)
	// 入口层解析出 seed 和代理，经 configs 透传给浏览器工厂。
	// seed 取值：环境变量 > 会话文件 > 新生成并写回，保证同一账号每次启动一致。
	configs.SetFingerprintSeed(configs.ResolveFingerprintSeed(
		cookies.NewLoadCookie(cookies.GetCookiesFilePath())))
	configs.SetProxy(configs.ProxyFromEnv())

	// 初始化服务
	xiaohongshuService := NewXiaohongshuService()

	// 创建并启动应用服务器
	appServer := NewAppServer(xiaohongshuService)

	switch transport {
	case "stdio":
		// 使用 stdio 传输模式
		if err := appServer.StartStdio(context.Background()); err != nil {
			logrus.Fatalf("MCP Server (stdio) error: %v", err)
		}
	case "http":
		// 使用 HTTP 传输模式
		if err := appServer.Start(port); err != nil {
			logrus.Fatalf("failed to run server: %v", err)
		}
	default:
		logrus.Fatalf("不支持的传输协议: %s, 请使用 http 或 stdio", transport)
	}
}

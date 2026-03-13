package main

import (
	"context"
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
)

func main() {
	var (
		headless  bool
		binPath   string // 浏览器二进制文件路径
		port      string
		transport string // 传输协议：http 或 stdio
	)
	flag.BoolVar(&headless, "headless", true, "是否无头模式")
	flag.StringVar(&binPath, "bin", "", "浏览器二进制文件路径")
	flag.StringVar(&port, "port", ":18060", "端口")
	flag.StringVar(&transport, "transport", "http", "传输协议: http 或 stdio")
	flag.Parse()

	if len(binPath) == 0 {
		binPath = os.Getenv("ROD_BROWSER_BIN")
	}

	configs.InitHeadless(headless)
	configs.SetBinPath(binPath)

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

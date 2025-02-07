package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
)

func main() {
	var parameters WhSvrParameters

	// 获取命令行参数
	flag.IntVar(&parameters.port, "port", 443, "Webhook server port.")                                                                        // 设置端口，默认 443
	flag.StringVar(&parameters.certFile, "tlsCertFile", "/etc/webhook/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")     // TLS 证书文件路径
	flag.StringVar(&parameters.keyFile, "tlsKeyFile", "/etc/webhook/certs/key.pem", "File containing the x509 private key to --tlsCertFile.") // TLS 私钥文件路径
	flag.Parse()

	// 加载 TLS 证书和私钥
	pair, err := tls.LoadX509KeyPair(parameters.certFile, parameters.keyFile)
	if err != nil {
		glog.Errorf("Failed to load key pair: %v", err) // 记录错误信息
		return
	}

	// 创建 webhook 服务器实例
	whsvr := &WebhookServer{
		server: &http.Server{
			Addr:      fmt.Sprintf(":%v", parameters.port),                // 绑定到指定的端口
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}}, // 配置 TLS 证书
		},
	}

	// 定义 HTTP 服务器和处理器
	mux := http.NewServeMux()                // 创建一个新的 ServeMux
	mux.HandleFunc("/mutate", whsvr.serve)   // 注册 "/mutate" 路由
	mux.HandleFunc("/validate", whsvr.serve) // 注册 "/validate" 路由
	whsvr.server.Handler = mux               // 设置 HTTP 服务器的 Handler

	// 在新 goroutine 中启动 webhook 服务器
	go func() {
		// 启动 HTTPS 服务
		if err := whsvr.server.ListenAndServeTLS("", ""); err != nil {
			glog.Errorf("Failed to listen and serve webhook server: %v", err) // 记录启动错误
		}
	}()

	glog.Info("Server started") // 记录服务器已启动

	// 监听操作系统关闭信号
	signalChan := make(chan os.Signal, 1)                      // 创建信号通道
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM) // 添加需要监听的信号
	<-signalChan                                               // 阻塞直到接收到信号

	glog.Infof("Got OS shutdown signal, shutting down webhook server gracefully...") // 日志记录关闭信号
	whsvr.server.Shutdown(context.Background())                                      // 优雅地关停 HTTP 服务器
}

package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zmall/ioc"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"zmall/config"
)

func main() {
	err := godotenv.Load("config/.env")
	if err != nil {
		panic(err)
	}

	app := ioc.NewApp()

	server := &http.Server{
		Addr:    config.GetConf().Server.Host,
		Handler: app.Server,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	zap.L().Info("Server is running", zap.String("address", "http://localhost:9000"))

	// 创建通道监听信号
	quit := make(chan os.Signal, 1)

	// 监听信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞直到收到信号
	<-quit
	zap.L().Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅地关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		zap.L().Error("Server forced shutting down", zap.Error(err))
	}

	zap.L().Info("Server exited gracefully")
}

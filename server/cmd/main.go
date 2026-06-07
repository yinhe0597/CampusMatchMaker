package main

import (
	"fmt"
	"os"

	v1 "campus_collab/api/v1"
	"campus_collab/internal/infra/config"
	"campus_collab/internal/infra/database"
	"campus_collab/internal/infra/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "配置加载失败: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化日志
	log := logger.InitLogger(cfg.Log)
	defer log.Sync()

	log.Info("校园协作平台启动中",
		zap.String("env", cfg.App.Env),
		zap.String("name", cfg.App.Name),
	)

	// 3. 初始化数据库
	db, err := database.InitMySQL(cfg.Database, log)
	if err != nil {
		log.Fatal("数据库初始化失败", zap.Error(err))
	}
	// 4. 初始化 Gin 引擎
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())

	// 5. 注册路由
	v1.RegisterRoutes(r, cfg, log, db)

	// 6. 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	log.Info("服务器启动", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		log.Fatal("服务器启动失败", zap.Error(err))
	}
}

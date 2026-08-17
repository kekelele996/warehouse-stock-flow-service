package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/wmsflow/wmsflow/internal/config"
	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/database"
	"github.com/wmsflow/wmsflow/internal/handler"
	"github.com/wmsflow/wmsflow/internal/middleware"
	"github.com/wmsflow/wmsflow/internal/repository"
	"github.com/wmsflow/wmsflow/internal/router"
	"github.com/wmsflow/wmsflow/internal/service"
	"github.com/wmsflow/wmsflow/internal/util"
)

func main() {
	logger := util.NewLogger(os.Getenv("LOG_LEVEL"))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config failed", "err", err)
		os.Exit(1)
	}

	db, err := database.Connect(cfg, logger)
	if err != nil {
		logger.Error("connect database failed", "err", err)
		os.Exit(1)
	}
	if err := database.Seed(context.Background(), db, logger); err != nil {
		logger.Error("seed data failed", "err", err)
		os.Exit(1)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Warn("redis unavailable, rate limit falls back to in-memory", "err", err)
		redisClient = nil
	} else {
		logger.Info(constants.LogRedisConnected, "host", cfg.RedisHost)
	}

	// 仓储
	userRepo := repository.NewUserRepository(db)
	ownerRepo := repository.NewOwnerRepository(db)
	productRepo := repository.NewProductRepository(db)
	binRepo := repository.NewBinLocationRepository(db)
	invRepo := repository.NewInventoryRepository(db)
	inboundRepo := repository.NewInboundRepository(db)
	outboundRepo := repository.NewOutboundRepository(db)
	logRepo := repository.NewOperationLogRepository(db)

	// 服务
	auditSvc := service.NewAuditService(logRepo, logger)
	inventorySvc := service.NewInventoryService(invRepo, binRepo, productRepo, ownerRepo, logger)
	authSvc := service.NewAuthService(db, cfg.JWTSecret, cfg.JWTExpire, userRepo, ownerRepo, auditSvc, logger)
	ownerSvc := service.NewOwnerService(ownerRepo, productRepo, inboundRepo, outboundRepo, invRepo, auditSvc, logger)
	productSvc := service.NewProductService(productRepo, ownerRepo, auditSvc, logger)
	binSvc := service.NewBinLocationService(db, binRepo, productRepo, invRepo, ownerRepo, auditSvc, logger)
	inboundSvc := service.NewInboundService(db, inboundRepo, ownerRepo, productRepo, binRepo, userRepo, inventorySvc, auditSvc, logger)
	outboundSvc := service.NewOutboundService(db, outboundRepo, ownerRepo, productRepo, binRepo, userRepo, inventorySvc, auditSvc, logger)
	dashSvc := service.NewDashboardService(inboundRepo, outboundRepo, binRepo, inventorySvc, logger)

	// 处理器
	authHandler := handler.NewAuthHandler(authSvc, logger)
	ownerHandler := handler.NewOwnerHandler(ownerSvc, logger)
	productHandler := handler.NewProductHandler(productSvc, logger)
	binHandler := handler.NewBinLocationHandler(binSvc, logger)
	invHandler := handler.NewInventoryHandler(inventorySvc, logger)
	inboundHandler := handler.NewInboundHandler(inboundSvc, logger)
	outboundHandler := handler.NewOutboundHandler(outboundSvc, logger)
	dashboardHandler := handler.NewDashboardHandler(dashSvc, logger)
	auditHandler := handler.NewAuditHandler(auditSvc, logger)

	// 中间件
	auditMW := middleware.Audit(logRepo, logger)
	rateLimitMW := middleware.RateLimit(redisClient, 120, time.Minute)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(rateLimitMW)
	router.Register(
		engine, cfg.JWTSecret, logger,
		authHandler, ownerHandler, productHandler, binHandler, invHandler,
		inboundHandler, outboundHandler, dashboardHandler, auditHandler, auditMW,
	)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: engine,
	}
	go func() {
		logger.Info(constants.LogServerStarted, "port", cfg.HTTPPort, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server listen failed", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(constants.LogServerShutdown, "err", err)
	}
}

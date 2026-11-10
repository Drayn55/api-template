package main

import (
	"os"

	"api-template/api"
	"api-template/configs"
	"api-template/internal/audit"
	"api-template/internal/handler"
	"api-template/internal/logger"
	"api-template/internal/repository"
	"api-template/internal/service"
	"api-template/internal/storage"
	"api-template/pkg/jwtutil"
	"api-template/pkg/migrationutil"
)

// @title API Template
// @version 1.0
// @description Standard enterprise Go Gin REST API starter template.
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Handle migration & seeder CLI commands (e.g. go run cmd/main.go migrate)
	migrationutil.HandleCLI()

	// Initialize structured logger
	logger.Init()

	// Load configuration
	cfg, err := configs.Load()
	if err != nil {
		logger.L().Error("failed to load configuration", "err", err.Error())
		os.Exit(1)
	}

	// Connect to Database
	db, err := configs.NewDB(cfg.DB)
	if err != nil {
		logger.L().Warn("failed to connect to database (continuing in degraded mode)", "err", err.Error())
	} else {
		defer db.Close()
		logger.L().Info("database connected successfully", "db", cfg.DB.DBName)
	}

	// Connect to MinIO (optional)
	var storageService *service.StorageService
	minioClient, err := configs.NewMinIO(cfg.MinIO)
	if err != nil {
		logger.L().Warn("failed to initialize minio client", "err", err.Error())
	} else if minioClient != nil {
		minioStorage := storage.NewMinioStorage(minioClient, cfg.MinIO.Bucket)
		storageService = service.NewStorageService(minioStorage)
		logger.L().Info("minio storage configured", "bucket", cfg.MinIO.Bucket)
	}

	// Initialize core utilities
	jwt := jwtutil.New(cfg.JWT.Secret)

	var auditService *audit.Service
	var apiLogRepo *repository.APILogRepo
	if db != nil {
		auditService = audit.NewService(audit.NewRepository(db))
		_ = auditService // ready to be passed into entity repositories
		apiLogRepo = repository.NewAPILogRepo(db)
	}

	// Initialize Handlers
	handlers := api.Handlers{
		Health: handler.NewHealthHandler(db),
		Upload: handler.NewUploadHandler(storageService),
	}

	// Initialize Router
	r := api.NewRouter(handlers, apiLogRepo, jwt)

	logger.L().Info("server running on", "addr", cfg.ServerAddr)
	if err := r.Run(cfg.ServerAddr); err != nil {
		logger.L().Error("failed to start server", "err", err.Error())
		os.Exit(1)
	}
}


package api

import (
	"api-template/internal/handler"
	"api-template/internal/middleware"
	"api-template/internal/repository"
	"api-template/pkg/jwtutil"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Health *handler.HealthHandler
	Upload *handler.UploadHandler
}

func NewRouter(h Handlers, apiLogRepo *repository.APILogRepo, jwt *jwtutil.JWTUtil) *gin.Engine {
	r := gin.New()

	// Base middlewares
	r.Use(middleware.CORS())
	_ = r.SetTrustedProxies(nil)

	r.Use(middleware.RequestLog())
	if apiLogRepo != nil {
		r.Use(middleware.APILogger(apiLogRepo))
	}
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.Recover())

	// Public health check endpoint
	if h.Health != nil {
		r.GET("/health", h.Health.Check)
	}

	// Protected API route group
	apiGroup := r.Group("/api")
	if jwt != nil {
		apiGroup.Use(middleware.Auth(jwt))
	}

	{
		// Storage / Upload endpoints
		if h.Upload != nil {
			uploadGroup := apiGroup.Group("/upload")
			{
				uploadGroup.POST("/presigned-url", h.Upload.GeneratePresignedURL)
			}
		}

		// Register your module routes here:
		// users := apiGroup.Group("/users")
		// {
		//     users.GET("", h.User.List)
		//     users.POST("", h.User.Create)
		// }
	}

	return r
}

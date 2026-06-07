package v1

import (
	"time"

	"campus_collab/internal/domain/classrepo"
	"campus_collab/internal/domain/poll"
	"campus_collab/internal/domain/timetable"
	"campus_collab/internal/domain/user"
	"campus_collab/internal/handler"
	"campus_collab/internal/handler/middleware"
	"campus_collab/internal/infra/config"
	"campus_collab/internal/service"
	"campus_collab/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RegisterRoutes 注册所有 v1 路由
func RegisterRoutes(r *gin.Engine, cfg *config.Config, log *zap.Logger, db *gorm.DB) {
	// 全局中间件
	r.Use(middleware.CORS(cfg.CORS.Origins))
	r.Use(RequestIDMiddleware())

	// 健康检查（不需要认证）
	r.GET("/api/v1/health", func(c *gin.Context) {
		response.OK(c, gin.H{
			"status":  "healthy",
			"version": "0.1.0",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// ===== 依赖注入 =====
	userRepo := user.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWT, cfg.Encrypt.Key)
	authHandler := handler.NewAuthHandler(authService)

	classRepo := classrepo.NewClassRepository(db)
	// 课表模块（需先创建以便注入 ClassService）
	ttRepo := timetable.NewTimetableRepository(db)
	ttService := service.NewTimetableService(ttRepo, classRepo)
	pollRepo := poll.NewPollRepository(db)
	pollService := service.NewPollService(pollRepo, classRepo, ttRepo)
	pollHandler := handler.NewPollHandler(pollService)

	classService := service.NewClassService(classRepo, ttService)
	classHandler := handler.NewClassHandler(classService)
	ttHandler := handler.NewTimetableHandler(ttService)

	// v1 API 路由组
	api := r.Group("/api/v1")
	{
		// 认证模块（无需 Token）
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 需要 JWT 认证的路由
		authorized := api.Group("")
		authorized.Use(middleware.JWTAuth(cfg.JWT))
		{
			authorized.GET("/auth/me", authHandler.GetCurrentUser)
			authorized.POST("/auth/refresh-token", authHandler.RefreshToken)

			// 班级模块路由
			classes := authorized.Group("/classes")
			{
				classes.GET("", classHandler.ListMyClasses)
				classes.POST("", classHandler.CreateClass)
				classes.GET("/by-code/:code", classHandler.LookupClassByCode)
				classes.GET("/:id", classHandler.GetClassDetail)
				classes.POST("/:id/join", classHandler.JoinClass)
				classes.GET("/:id/members", classHandler.ListMembers)
				classes.DELETE("/:id/members/:userId", classHandler.RemoveMember)
			}

			// 课表模块路由
			timetables := authorized.Group("/timetables")
			{
				timetables.POST("/class/:classId", ttHandler.CreateClassTimetable)
				timetables.GET("/class/:classId", ttHandler.GetClassTimetable)
				timetables.PUT("/class/:classId", ttHandler.UpdateClassTimetable)
				timetables.POST("/personal", ttHandler.CreatePersonalTimetable)
				timetables.GET("/personal", ttHandler.GetPersonalTimetable)
				timetables.PUT("/personal/:id", ttHandler.UpdatePersonalTimetable)
				timetables.DELETE("/personal/:id", ttHandler.DeletePersonalTimetable)
				timetables.POST("/corrections", ttHandler.CreateCorrection)
				timetables.GET("/corrections", ttHandler.ListCorrections)
				timetables.PUT("/corrections/:id", ttHandler.ReviewCorrection)
			}

			// 投票模块路由
			polls := authorized.Group("/polls")
			{
				polls.POST("", pollHandler.CreatePoll)
				polls.GET("", pollHandler.ListPolls)
				polls.GET("/:id", pollHandler.GetPollDetail)
				polls.PUT("/:id", pollHandler.EditPoll)
				polls.POST("/:id/open", pollHandler.OpenPoll)
				polls.POST("/:id/close", pollHandler.ClosePoll)
				polls.GET("/:id/options", pollHandler.GetOptions)
				polls.POST("/:id/vote", pollHandler.SubmitVote)
				polls.GET("/:id/results", pollHandler.GetResults)
				polls.POST("/:id/finalize", pollHandler.FinalizePoll)
			}
		}
	}

	log.Info("路由注册完成", zap.Int("routes_count", len(r.Routes())))
}

// RequestIDMiddleware 请求追踪 ID 中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func generateRequestID() string {
	return time.Now().Format("20060102150405.000")
}

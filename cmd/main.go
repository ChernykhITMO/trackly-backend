package main

import (
	"flag"
	"github.com/labstack/echo/v4"
	"log"
	"trackly-backend/internal/api"
	"trackly-backend/internal/config"
	"trackly-backend/internal/db"
	"trackly-backend/internal/middleware"
	"trackly-backend/internal/repositories"
)

type Server struct {
	*api.UserApi
	*api.HabitsApi
	*api.StatisticApi
	*api.ProgressApi
}

func main() {
	// Загрузка конфигурации
	configFilePath := flag.String("configs", "", "Path to the configuration file")
	flag.Parse()
	println("config path:" + *configFilePath)

	cfg, err := config.LoadConfig(*configFilePath)
	if err != nil {
		log.Fatalf("Could read config: %v", err)
	}

	// Инициализация базы данных
	database, err := db.InitDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Инициализация Echo
	e := echo.New()

	// Инициализация репозитория и сервера
	userRepo := repositories.NewUserRepository(database)
	minio, err := db.NewMinioClient(cfg)
	if err != nil {
		log.Fatalf("Could not connect to S3: %v", err)
	}

	userApi := api.NewUserApi(userRepo, cfg, minio)

	planRepo := repositories.NewPlanRepository(database)

	habitRepo := repositories.NewHabitRepository(database)
	habitsApi := api.NewHabitsApi(habitRepo, planRepo)

	statisitcRepo := repositories.NewStatisticRepository(database)

	statisticApi := api.NewStatisticApi(habitRepo, statisitcRepo)
	progressApi := api.NewProgressApi(statisitcRepo, habitRepo, planRepo)
	server := &Server{userApi, habitsApi, statisticApi, progressApi}

	RegisterHandlers(e, server, cfg.JwtSecret)

	// Запуск сервера
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}

func RegisterHandlers(router *echo.Echo, si api.ServerInterface, jwtSecret string) {
	// Public routes (no authentication required)
	wrapper := api.ServerInterfaceWrapper{
		Handler: si,
	}
	router.Use(middleware.Cors())

	publicGroup := router.Group("")
	publicGroup.POST("/api/auth/login", wrapper.PostApiAuthLogin)
	publicGroup.POST("/api/auth/register", wrapper.PostApiAuthRegister)

	// Protected routes (JWT authentication required)
	protectedGroup := router.Group("")

	protectedGroup.Use(middleware.JWTMiddleware([]byte(jwtSecret))) // Apply JWT middleware to protected routes

	protectedGroup.GET("/api/habits", wrapper.GetApiHabits)
	protectedGroup.POST("/api/habits", wrapper.PostApiHabits)
	protectedGroup.GET("/api/habits/:habitId", wrapper.GetApiHabitsHabitId)
	protectedGroup.PUT("/api/habits/:habitId", wrapper.PutApiHabitsHabitId)
	protectedGroup.POST("/api/habits/:habitId/score", wrapper.PostApiHabitsHabitIdScore)
	protectedGroup.GET("/api/habits/:habitId/statistic", wrapper.GetApiHabitsHabitIdStatistic)
	protectedGroup.GET("/api/habits/:habitId/statistic/total", wrapper.GetApiHabitsHabitIdStatisticTotal)
	protectedGroup.POST("/api/users/avatar", wrapper.PostApiUsersAvatar)
	protectedGroup.GET("/api/users/profile", wrapper.GetApiUsersProfile)
	protectedGroup.PUT("/api/users/profile", wrapper.PutApiUsersProfile)

	protectedGroup.POST("/api/users/avatar", wrapper.PostApiUsersAvatar)
	protectedGroup.GET("/api/users/avatar", wrapper.GetApiUsersAvatar)

}

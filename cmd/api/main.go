package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/irpanzy/Task-Forge/internal/config"
	"github.com/irpanzy/Task-Forge/internal/controller"
	"github.com/irpanzy/Task-Forge/internal/middleware"
	"github.com/irpanzy/Task-Forge/internal/repository"
	"github.com/irpanzy/Task-Forge/internal/route"
	"github.com/irpanzy/Task-Forge/internal/service"
	"github.com/irpanzy/Task-Forge/pkg/response"
)

func main() {
	// 1. Load config & connect to Neon PostgreSQL
	config.LoadEnv()
	config.ConnectDB()

	// 2. Inisialisasi Fiber app
	app := fiber.New()

	// 3. Global Middlewares
	app.Use(logger.New())
	app.Use(middleware.NewCORS())
	app.Use(middleware.NewCSRF())

	app.Get("/", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, "TaskForge API is running", fiber.Map{
			"app":     "TaskForge API",
			"version": "1.0.0",
		})
	})

	// 4. Dependency Injection (Wiring Layer)
	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	authController := controller.NewAuthController(userService)
	userController := controller.NewUserController(userService)

	boardRepo := repository.NewBoardRepository(config.DB)
	boardMemberRepo := repository.NewBoardMemberRepository(config.DB)
	boardService := service.NewBoardService(boardRepo, boardMemberRepo, userRepo)
	boardController := controller.NewBoardController(boardService)

	// 5. Setup Routes
	route.SetupRoutes(app, authController, userController, boardController)

	// 6. Start Server
	port := config.AppConfig.Port
	if port == "" {
		port = "3000"
	}

	log.Printf("Server TaskForge berjalan di http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}

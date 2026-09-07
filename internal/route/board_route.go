package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/irpanzy/Task-Forge/internal/controller"
	"github.com/irpanzy/Task-Forge/internal/middleware"
)

func RegisterBoardRoutes(router fiber.Router, boardCtrl controller.BoardController) {
	boards := router.Group("/boards", middleware.Authenticate())

	boards.Post("/", boardCtrl.CreateBoard)
	boards.Get("/", boardCtrl.GetBoards)
	boards.Get("/:id", boardCtrl.GetBoardByID)
	boards.Put("/:id", boardCtrl.UpdateBoard)
	boards.Delete("/:id", boardCtrl.DeleteBoard)

	// Board Members
	boards.Post("/:id/members", boardCtrl.AddMembers)
	boards.Get("/:id/members", boardCtrl.GetMembers)
	boards.Delete("/:id/members/:userId", boardCtrl.DeleteMember)
}

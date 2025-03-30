package routes

import (
	"chat/internal/database"
	"chat/internal/models"
	"github.com/gofiber/fiber/v2"
)

func ProtectedRoutes(app fiber.Router) {
	database.InitDB()
	app.Post("/createpost", CreatePost)
	app.Get("/getpost", GetPost)
	app.Post("/like", Like)
	app.Post("/unlike", Unlike) // Добавили маршрут
}

func CreatePost(c *fiber.Ctx) error {

	var post models.Post
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	database.DB.Create(&post)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"post": post})
}

func GetMyPost(c *fiber.Ctx) error {
	var posts []models.Post
	if err := database.DB.Where("user_id = ?", c.Params("user_id")).First(&posts); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"posts": posts})
}

func GetPost(c *fiber.Ctx) error {
	var posts []models.Post

	if err := database.DB.Find(&posts).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"posts": posts})
}

func Like(c *fiber.Ctx) error {
	var like models.Like
	if err := c.BodyParser(&like); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var post models.Post
	if err := database.DB.First(&post, like.PostID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}

	database.DB.Create(&like)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Post liked", "like": like})
}

func Unlike(c *fiber.Ctx) error {
	var like models.Like
	if err := c.BodyParser(&like); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := database.DB.Where("post_id = ? AND user_id = ?", like.PostID, like.UserID).Delete(&models.Like{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to unlike"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Like removed"})
}

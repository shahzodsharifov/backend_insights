package controllers

import (

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wpcodevo/golang-fiber/initializers"
	"github.com/wpcodevo/golang-fiber/models"
)

func CreateComment(c *fiber.Ctx) error {

	postID := c.Params("postID")

	
	var comment models.Comment
	if err := c.BodyParser(&comment); err !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	
	if comment.Body == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Comment body is required",
		})
	}
	postIDParsed, parseErr := uuid.Parse(postID)
	if parseErr !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid post id",
		})
	}

	authorIDParsed, authorParseErr := uuid.Parse(comment.AuthorID.String())
	if authorParseErr !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid author id",
		})
	}
	comment.PostID = postIDParsed 
	comment.AuthorID = authorIDParsed

	if err := initializers.DB.Create(&comment).Error; err !=nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create comment",
		})
	}

	return c.JSON(comment)

}

func GetPostComments(c *fiber.Ctx) error {

	postID := c.Params("postID")
	var post models.Post
	// var comment models.Post
	if err := initializers.DB.Where("id = ? ", postID).Preload("Comments").First(&post).Error; err !=nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Post not found",

		})
	}

	return c.JSON(post.Comments)
}
package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wpcodevo/golang-fiber/initializers"
	"github.com/wpcodevo/golang-fiber/models"
	"gorm.io/gorm"
)

func AddPost(c *fiber.Ctx) error {
	userID := c.Params("userID")

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err !=nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving user",
		})
	} 

	var post models.Post
	if err := c.BodyParser(&post); err !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing post data",
		})
	}

	post.AuthorID = user.ID

	if err := initializers.DB.Create(&post).Error; err !=nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating post",
		})
	}

	user.Posts = append(user.Posts, post)

	return c.JSON(fiber.Map{"post": post})
}

func GetUserPosts(c *fiber.Ctx) error {

	userID := c.Params("userID")

	var user models.User
	if err := initializers.DB.Where("id = ?", userID).Preload("Posts").First(&user).Error; err !=nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})

	}

	return c.JSON(user.Posts)
}

func GetPostByID(c *fiber.Ctx) error {
	postID := c.Params("postID")

	var post models.Post
	res := initializers.DB.First(&post, "id =?", postID)
	if err := res.Error; err !=nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": err.Error()})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"post": post}})
}

func GetTrendingPosts(c *fiber.Ctx) error {
	var posts *[]models.Post

	res := initializers.DB.Limit(10).Find(&posts)
	if err := res.Error; err !=nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"post": posts}})
}

func GetCategoryPosts(c *fiber.Ctx) error {

	category := c.Params("category")
	// var payload *models.CategoryInput

	// if err := c.BodyParser(&payload); err !=nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	// }
	var posts *[]models.Post

	res := initializers.DB.Find(&posts, "category=?", category)
	if err := res.Error; err !=nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": err.Error()})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"post": posts}})

}

func LikePost(c *fiber.Ctx) error {
	// likedPostID := c.Params("likedPostID")
	// userID := c.Params("userID")
	var payload *models.Like
	var like models.Like
	if err := c.BodyParser(&payload); err !=nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	likedPostIDParsed, parseErr := uuid.Parse(payload.PostID.String())
	if parseErr !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid post id",
		})
	}

	userIDParsed, userIDParsedErr := uuid.Parse(payload.UserID.String())
	if userIDParsedErr !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user id",
		})
	}

	like.PostID = likedPostIDParsed
	like.UserID = userIDParsed

	if err := initializers.DB.Create(&like).Error; err !=nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create like",
		})

		
	}
	return c.JSON(like)
}

func Unlike(c *fiber.Ctx) error {
	var payload *models.Like
	var like models.Like
	

	if err := c.BodyParser(&payload); err !=nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := initializers.DB.Where("user_id = ? AND post_id = ?", payload.UserID, payload.PostID).Delete(&like).Error; err !=nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete like",
		})

	
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Deleted it successfully",
	})

	
}

func GetLikes(c*fiber.Ctx) error {
	likedPostID := c.Params("postID")
	var post models.Post

	if err := initializers.DB.Where("id = ?", likedPostID).Preload("Likes").First(&post).Error; err !=nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Post not found",
		})
	}

	return c.JSON(post.Likes)
}


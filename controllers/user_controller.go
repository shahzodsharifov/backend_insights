package controllers

import (
	"fmt"

	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/wpcodevo/golang-fiber/initializers"
	"github.com/wpcodevo/golang-fiber/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUpUser(c *fiber.Ctx) error {
	var payload *models.SignUpInput

	if err := c.BodyParser(&payload); err !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	fmt.Println(payload)
	now := time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	newUser:= models.User{
		Name: payload.Name,
		Username: strings.ToLower(payload.Username),
		Email: strings.ToLower(payload.Email),
		Password: string(hashedPassword),
		Photo: &payload.Photo,
		Role: &payload.Role,
		CreatedAt: &now,
		
	}

	result := initializers.DB.Create(&newUser)

	if result.Error !=nil && strings.Contains(result.Error.Error(), "duplicate key violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": newUser}})
}


func SignInUser(c *fiber.Ctx) error {
	var payload *models.SignInInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	var user models.User
	result := initializers.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid email or Password"})

	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})

	}

	config, _ := initializers.LoadConfig(".")

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)


	claims["sub"] = user.ID
	claims["exp"] = now.Add(config.JwtExpiresIn).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(config.JwtSecret))

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})

	}

	c.Cookie(&fiber.Cookie{
		Name: "token",
		Value: tokenString,
		Path: "/",
		MaxAge: config.JwtMaxAge * 60,
		Secure: true,
		HTTPOnly: true,
		Domain: "insights-frontend-55ak.onrender.com/",
		SameSite: "none",

	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "token": tokenString})
}

func LogoutUser(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name: "token",
		Value: "",
		Expires: expired,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}


func GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserResponse)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}


func FindAllUsers(c*fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _:= strconv.Atoi(page)
	intLimit, _:= strconv.Atoi(limit)
	offset := (intPage -1) * intLimit

	var users []models.User
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&users)
	if results.Error !=nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(users), "users": users})


}


func UpdateUser(c *fiber.Ctx) error {
	userId := c.Params("userId")

	var payload  *models.SignUpInput

	if err:= c.BodyParser(&payload); err !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var user models.User
	result := initializers.DB.First(&user, "id = ?", userId)
	if err := result.Error; err !=nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": err.Error()})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.Name !="" {
		updates["name"] = payload.Name
	}
	if payload.Username !="" {
		updates["username"] = payload.Username
	}
	if payload.Email !="" {
		updates["email"] = payload.Email
	}

	initializers.DB.Model(&user).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}


func FindUserById(c *fiber.Ctx) error {
	userId := c.Params("userId")


	
	var user models.User
	result := initializers.DB.First(&user, "id = ?", userId)
	if err := result.Error; err !=nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": err.Error()})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}

func FindUserByUsername(c *fiber.Ctx) error {
	userName :=c.Params("userName")
	ser := "%"+userName+"%"
	var users *[]models.User
	res := initializers.DB.Limit(3).Where("username LIKE ?",ser).Find(&users)
	if err :=res.Error; err !=nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": err.Error()})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"users": users}})


}

func GetUserFollowers(c *fiber.Ctx) error {
	userId :=c.Params("userId")



	var rels *[]models.UserRelationship
	
	// if err := initializers.DB.Preload("Followers").Where("id = ?", payload.UserID).First(&user).Error; err !=nil {
	// 	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	// }

	// return c.JSON(user.Followers)

	resErr := initializers.DB.Find(&rels, "following_id = ?", userId);
	if resErr.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": resErr.Error})
	}
	return c.JSON(rels)
}

func GetUserFollowing(c *fiber.Ctx) error {
	userId :=c.Params("userId")

	// var payload models.FindUserInput
	var rels *[]models.UserRelationship
	// if err:= c.BodyParser(&payload); err !=nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	// }

	// var user models.User
	// if err := initializers.DB.Preload("Following").Where("id = ?", payload.UserID).First(&user).Error; err !=nil {
	// 	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})

	// }
	
	resErr := initializers.DB.Find(&rels, "follower_id = ?", userId);
	if resErr.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": resErr.Error})
	}
	return c.JSON(&rels)

}

func AddFollower(c *fiber.Ctx) error {
	// followerID  := c.Params("followerID")
	// followingID := c.Params("followingID")

	var payload *models.FollowInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	followerUUID, err := uuid.Parse(payload.FollowingID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid follower ID"})
	}

	followingUUID, err := uuid.Parse(payload.FollowingID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid following ID"})
	}
	var follower, following models.User

	if err := initializers.DB.Where("id = ?", followerUUID).First(&follower).Error ; err !=nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Follower not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add follower"})
	}

	if err := initializers.DB.Where("id = ?", followingUUID).First(&following).Error ; err !=nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Follower not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add following"})

	}

	var relationship models.UserRelationship

	// //check if the relationship exists
	// if err := initializers.DB.Where("follower_id = ? AND following_id = ?", payload.FollowerID, payload.FollowingID ).First(&relationship).Error; err == nil {
	// 	return c.Status(fiber.StatusConflict).JSON(fiber.Map{
	// 		"message": "Relationship already exists"})
	// }

	//Create the relationship
	relationship = models.UserRelationship{
		ID: uuid.New(),
		FollowerID: followerUUID,
		FollowingID: followingUUID,
	}
	//Add the follower to following user's followers
	initializers.DB.Model(&following).Association("Followers").Append(&follower)


	if err := initializers.DB.Create(&relationship).Error; err !=nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add follower"})

	}
	initializers.DB.Create(&relationship)

	


	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Follower added successfully"})
}

func AddFollowing(c *fiber.Ctx) error {
	// followerID := c.Params("followerID")
	// followingID := c.Params("followingID")

	var payload *models.FollowInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	var follower, following models.User	

	// followerUUID, err := uuid.Parse(followerID)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid follower ID"})
	// }

	// followingUUID, err := uuid.Parse(followingID)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid following ID"})
	// }

	if err := initializers.DB.Where("id = ?", payload.FollowerID).First(&follower).Error ; err !=nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Follower not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add follower"})
	}

	if err := initializers.DB.Where("id = ?", payload.FollowingID).First(&following).Error ; err !=nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Following not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add following"})

	}



		var relationship models.UserRelationship
		followerIDParsed, parseErr := uuid.Parse(payload.FollowerID.String())
		followingIDParsed, parseErr := uuid.Parse(payload.FollowingID.String())

		if parseErr !=nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "could not parse the input",
			})
		}
		// //check if the relationship exists
		if err := initializers.DB.Where("follower_id = ? AND following_id = ?", payload.FollowerID, payload.FollowingID ).First(&relationship).Error; err == nil {
		if err == gorm.ErrRecordNotFound {
						//Create the relationship
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"message": "Relationship does not exist",
				})


			}
			
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Relationship already exists"})
		}
	
		//Create the relationship

		relationship = models.UserRelationship{
			ID: uuid.New(),
			FollowerID: followerIDParsed,
			FollowingID: followingIDParsed,
		}
		 dataErr := initializers.DB.Create(&relationship).Error; if dataErr !=nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not create a relationship"})
		}
	


	//Add the following to follower user's following
	// initializers.DB.Model(&follower).Association("Following").Append(&following)



	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Follower added successfully!!!"})


		


		
}

func DeleteFollowing(c *fiber.Ctx) error {
	var payload *models.FollowInput

	if err := c.BodyParser(&payload); err !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})

	}

	errors := models.ValidateStruct(payload)
	if errors !=nil {
	return c.Status(fiber.StatusBadRequest).JSON(errors)
}

	var follower, following models.User	
	
	if err := initializers.DB.Where("id = ?", payload.FollowerID).First(&follower).Error ; err !=nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Follower not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add follower"})
	}

	if err := initializers.DB.Where("id = ?", payload.FollowingID).First(&following).Error ; err !=nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Following not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add following"})

	}

	var relationship models.UserRelationship
	followerIDParsed, parseErr := uuid.Parse(payload.FollowerID.String())
	followingIDParsed, parseErr := uuid.Parse(payload.FollowingID.String())


	if parseErr !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "could not parse the input",
		})
	}

	res := initializers.DB.Delete(&relationship, "follower_id = ? AND following_id = ?", followerIDParsed, followingIDParsed)
	if res.Error !=nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not delete a relationship"})

	}

	return c.JSON("relationship deleted successfully")


}


	 	
func GetAllCompanies(c*fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _:= strconv.Atoi(page)
	intLimit, _:= strconv.Atoi(limit)
	offset := (intPage -1) * intLimit

	var companies []models.User
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&companies, "role = ?", "company")
	if results.Error !=nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error })
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(companies), "companies":companies})
}


//0030db83-1a67-4f78-846b-76a01a4c49c0
//83a8d136-607e-4a72-8d86-30c9e5d83b9a


package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/wpcodevo/golang-fiber/controllers"
	"github.com/wpcodevo/golang-fiber/initializers"
	"github.com/wpcodevo/golang-fiber/middleware"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())

	}
	initializers.ConnectDB(&config)
}

func main() {
	app := fiber.New()
	micro := fiber.New()

	app.Mount("/api", micro)
	app.Use(logger.New())

	micro.Use(cors.New(cors.Config{
		AllowOrigins:     "https://insights-frontend-55ak.onrender.com",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://insights-frontend-55ak.onrender.com",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	micro.Get("/api/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Welcome to Golang",
		})
	})

	micro.Route("/auth", func(router fiber.Router) {
		router.Post("/register", controllers.SignUpUser)
		router.Post("/login", controllers.SignInUser)
		router.Get("/logout", middleware.DeserializeUser)
		router.Post("/follow", controllers.AddFollowing)

	})

	micro.Get("/users/me", middleware.DeserializeUser, controllers.GetMe)

	micro.Route("/users", func(router fiber.Router) {
		router.Post("/", controllers.SignUpUser)
		router.Get("/", controllers.FindAllUsers)
		

	})


	micro.Route("/posts", func (router fiber.Router)  {
		router.Get("/trendingPosts", controllers.GetTrendingPosts)
		router.Get("/topics/:category", controllers.GetCategoryPosts)
		router.Get("/:postID", controllers.GetPostByID)
		router.Post("/:postID/likePost", controllers.LikePost)
		router.Post("/:postID/unlikePost", controllers.Unlike)
		router.Get("/:postID/likes", controllers.GetLikes)
		router.Get("/:postID/comments", controllers.GetPostComments)
		router.Post("/:postID/createComment", controllers.CreateComment)

	})

	micro.Route("/companies", func (router fiber.Router) {
		router.Get("/", controllers.GetAllCompanies)
})




	micro.Route("/vaccancies", func (router fiber.Router) {
		router.Get("/", controllers.GetAllVacancies)
		router.Get("/:vaccancyID", controllers.GetVaccancyByID)
	})

	micro.Route("/events", func (router fiber.Router) {
		router.Get("/", controllers.GetAllEvents)
		router.Get("/:eventID", controllers.GetEventsByID)
	})

	micro.Route("/users/search/:userName", func(router fiber.Router) {
		router.Get("/", controllers.FindUserByUsername)
	}) 
	micro.Route("/users/:userId", func(router fiber.Router) {
		router.Get("", controllers.FindUserById)
		router.Patch("", controllers.UpdateUser)
		router.Post("/follow", controllers.AddFollower)
		router.Post("/following", controllers.AddFollowing)

		router.Get("/followers/", controllers.GetUserFollowers)
		router.Get("/following/", controllers.GetUserFollowing)
		router.Post("/unfollow/", controllers.DeleteFollowing)
		
		router.Post("/addPost/", controllers.AddPost)
		router.Get("/posts", controllers.GetUserPosts)
		

		router.Get("/posts/:postID/comments", controllers.GetPostComments)
		router.Post("/posts/:postID/createComment", controllers.CreateComment)

		router.Post("/addVaccancy/", controllers.AddVaccancy)
		router.Get("/vaccancies", controllers.GetUserVaccancies)
		router.Get("/vaccancies/:vaccancyID", controllers.GetVaccancyByID)

		router.Post("/addEvent/", controllers.AddEvent)
		router.Get("/events", controllers.GetUserEvents)
		router.Get("/events/:eventID", controllers.GetEventsByID)

	})

	log.Fatal(app.Listen(":8000"))
}

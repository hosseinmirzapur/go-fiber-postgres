package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/hosseinmirzapur/simple-rest-api/models"
	"github.com/hosseinmirzapur/simple-rest-api/storage"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type Book struct {
	Author    string `json:"author" validation:"required"`
	Title     string `json:"title" validation:"required"`
	Publisher string `json:"publisher" validation:"required"`
}

func (r *Repository) CreateBook(c *fiber.Ctx) error {
	book := Book{}

	err := c.BodyParser(&book)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"messsage": "data not valid",
		})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "error creating book",
		})
		return err
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "the book has been added",
	})
}

func (r *Repository) GetBooks(c *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "error getting books",
		})
		return err
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"data": bookModels,
	})
}

func (r *Repository) GetBookById(c *fiber.Ctx) error {
	bookModel := &models.Books{}
	id := c.Params("id")
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id is required",
		})
		return nil
	}
	err := r.DB.Find(bookModel, id).Error

	if err != nil {
		c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "no book with such id found",
		})
		return err
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"data": bookModel,
	})
}

func (r *Repository) DeleteBook(c *fiber.Ctx) error {
	bookModel := &models.Books{}
	id := c.Params("id")
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id is required",
		})
		return nil
	}

	err := r.DB.Delete(bookModel, id).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "error deleting book",
		})
		return err
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "the book has been deleted",
	})
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/create-books", r.CreateBook)
	api.Delete("/delete-books/:id", r.DeleteBook)
	api.Get("/get-books/:id", r.GetBookById)
	api.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	db, err := storage.NewConnection(&config)

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	err = models.MigrateBooks(db)

	if err != nil {
		log.Fatal("Error migrating books")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()

	r.SetupRoutes(app)

	err = app.Listen(":3000")

	if err != nil {
		log.Fatal("Error starting server")
	}
}

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("hello world")
	app := fiber.New()
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	todos := []*Todo{}
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := Todo{}
		if err := c.BodyParser(&todo); err != nil {
			return err
		}
		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
		}
		todo.ID = len(todos) + 1
		todos = append(todos, &todo)
		return c.Status(200).JSON(todo)
	})
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return err
		}
		for _, todo := range todos {
			if todo.ID == id {
				todo.Completed = !todo.Completed
				//todos[i] = todo
				return c.Status(200).JSON(todos)
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return err
		}
		for i, todo := range todos {
			if todo.ID == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(todos)
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}

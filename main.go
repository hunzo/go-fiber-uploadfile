package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, //limit 10MB
	})
	app.Use(logger.New())

	app.Static("/static", "./uploads")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"info": "upload example",
		})
	})

	app.Get("/listfiles", func(c *fiber.Ctx) error {
		dirName := "./uploads"

		f, err := os.Open(dirName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		files, err := f.ReadDir(-1)
		f.Close()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		fileList := []string{}

		for _, file := range files {
			fileList = append(fileList, file.Name())
			// fmt.Println(file.Name())
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": fileList,
		})
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		uniqueID := uuid.New()

		fileName := strings.Replace(uniqueID.String(), "-", "", -1)
		fileExtension := strings.Split(file.Filename, ".")[1]
		_File := fmt.Sprintf("%s.%s", fileName, fileExtension)

		filePath := "./uploads/" + _File

		err = c.SaveFile(file, filePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		fileURL := fmt.Sprintf("http://localhost:8080/static/%s", _File)

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success":          true,
			"file_name":        file.Filename,
			"file_name_encode": _File,
			"file_sizes":       file.Size,
			"file_header":      file.Header,
			"file_url":         fileURL,
		})

	})

	app.Delete("/delete/:filename", func(c *fiber.Ctx) error {
		filename := c.Params("filename")

		err := os.Remove(fmt.Sprintf("./uploads/%s", filename))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success":  false,
				"messages": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"message": "file: " + filename + " has been deleted.",
		})
	})

	app.Listen(":8080")
}

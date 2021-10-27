package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

func main() {

	SERVER_NAME := os.Getenv("SERVER_NAME") //http://localhost
	SERVER_PORT := os.Getenv("SERVER_PORT") //8080
	SERVER := SERVER_NAME + ":" + SERVER_PORT

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, //limit 10MB
	})
	app.Use(logger.New())

	app.Static("/static", "./uploads")
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"info":     "upload example",
			"server":   SERVER,
			"listfile": "/listfiles",
			"upload":   "POST /upload",
			"delete":   "DELETE /delete/:filename",
			"clear":    "GET /clear",
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

		fileURL := fmt.Sprintf("%s/static/%s", SERVER, _File)

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

	app.Get("/clear", func(c *fiber.Ctx) error {
		directry := "./uploads/"

		dirRead, err := os.Open(directry)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success":  false,
				"messages": err.Error(),
			})
		}

		files, err := dirRead.ReadDir(0)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success":  false,
				"messages": err.Error(),
			})
		}

		var files_remove []string

		for index := range files {
			file := files[index]
			fileName := file.Name()
			filePath := directry + fileName

			os.Remove(filePath)
			files_remove = append(files_remove, filePath)
			fmt.Printf("remove file: %s\n", filePath)
		}

		return c.JSON(fiber.Map{
			"success":      true,
			"message":      "clear file",
			"remove_files": files_remove,
		})
	})

	app.Listen(":8080")
}

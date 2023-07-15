package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	img "image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	file, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.Println(err.Error())
	}
	logger.SetOutput(file)

	engine := html.New("./", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
		// the second parameter is data send to view
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		var Input struct {
			ImageName string
		}
		if err := c.BodyParser(&Input); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		image, err := c.FormFile("image")
		// from ajax in index.html
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		fmt.Printf("Image Name: %s\n", image.Filename)
		fmt.Printf("Image Size (bytes): %d\n", image.Size)
		fmt.Printf("Image Size (kB): %d\n", image.Size/1024)
		fmt.Printf("Image Size (mB): %f\n", (float64(image.Size)/1024)/1024)
		fmt.Printf("Mime Type: %s\n", image.Header.Get("Content-Type"))

		splitFileName := strings.Split(image.Filename, ".")
		extension := splitFileName[len(splitFileName)-1]
		fmt.Println(extension)

		newFileName := fmt.Sprintf("%s.%s", time.Now().Format("2006-01-02-15-04-05"), extension)
		fmt.Println(newFileName)

		fileHeader, _ := image.Open()
		defer fileHeader.Close()

		imageDimension, _, err := img.DecodeConfig(fileHeader)
		if err != nil {
			logger.Println(err)
		}
		width := imageDimension.Width
		height := imageDimension.Height
		fmt.Printf("width %d\n", width)
		fmt.Printf("height %d\n", height)

		// make folder in root dir
		folderUpload := filepath.Join(".", "uploads")
		if err := os.MkdirAll(folderUpload, 0770); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		// save img to uploads dir
		if err := c.SaveFile(image, "./uploads/"+newFileName); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"title":      Input.ImageName,
			"image_name": newFileName,
			"message":    "Successfully Upload Image",
		})
	})

	logger.Println("Server Running on Port 8080")
	app.Listen(":8080")
}

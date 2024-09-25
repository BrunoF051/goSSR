package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

var huggingFaceKey string

func main() {

	// Load environment variables from
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	huggingFaceKey = os.Getenv("HUGGINGFACE_KEY")
	if huggingFaceKey == "" {
		log.Fatal("HUGGINGFACE_KEY is not set in the environment")
	}

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", handleIndex)
	app.Post("/upload", handleUpload)

	log.Fatal(app.Listen(":3000"))
}

func handleIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":       "Template Example",
		"Description": "An example template",
		"Greeting":    "subtitle",
		"ResetForm":   false,
		"Error":       "",
	})
}

func handleUpload(c *fiber.Ctx) error {
	// Get the file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Render("index", fiber.Map{
			"Error": "Failed to get file from form",
		})
	}

	// Get the name from form
	name := c.FormValue("name")

	// Validate file type
	if !isValidImageType(file.Filename) {
		return c.Render("index", fiber.Map{
			"Error": "Invalid file type. Only JPG, JPEG, and PNG are allowed",
		})
	}

	// Validate file size (5MB max)
	if file.Size > 5*1024*1024 {
		return c.Render("index", fiber.Map{
			"Error": "File size exceeds 5MB limit",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Failed to open file",
		})
	}
	defer src.Close()

	// Read the file content
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Failed to read file",
		})
	}

	// Encode to base64
	base64String := base64.StdEncoding.EncodeToString(fileBytes)

	// Send to Hugging Face API
	processedImage, err := sendToHuggingFace(base64String)
	if err != nil {
		return c.Render("index", fiber.Map{
			"Error":     "Failed to process image: " + err.Error(),
			"ResetForm": true,
		})
	}

	// The processedImage is now a []byte containing the processed image data
	// You might want to save this or convert it back to base64 for display

	processedBase64 := base64.StdEncoding.EncodeToString(processedImage)

	// Render the index page with a success message
	return c.Render("index", fiber.Map{
		"Title":          "Template Example",
		"Description":    "An example template",
		"Greeting":       "subtitle",
		"Success":        "File uploaded and processed successfully",
		"UploadedName":   name,
		"FileName":       base64String,
		"ProcessedImage": processedBase64,
		"ResetForm":      true,
	})

}

func isValidImageType(filename string) bool {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg", ".png":
		return true
	default:
		return false
	}
}

func sendToHuggingFace(base64Image string) ([]byte, error) {
	url := "https://api-inference.huggingface.co/models/timbrooks/instruct-pix2pix"

	// Prepare the request payload
	payload := map[string]interface{}{
		"inputs": map[string]interface{}{
			"image":  base64Image,
			"prompt": "Make this image look like a professional profile picture with soft lighting and neutral background.",
			"parameters": map[string]interface{}{
				"guidance_scale":      7.5,
				"num_inference_steps": 10,
				"width":               512,
				"height":              512,
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Create the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer huggingFaceKey")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Use-Cache", "true")
	req.Header.Set("X-Wait-For-Model", "true")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

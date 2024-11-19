package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type GatewayController struct{}

func (gc *GatewayController) Index(c *fiber.Ctx) error {
	method := c.Method()
	path := c.Path()
	body := c.Body()
	ip := c.IP()

	log.Printf("Received Request :[ %s ] %s %s", ip, method, path)

	client := &http.Client{
		Timeout: 500 * time.Second,
	}
	req, err := http.NewRequest(
		method,
		"https://streambe01.indonesia.us.kg"+path,
		bytes.NewReader(body),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "ERROR_GENERATING_REQUEST",
		})
	}

	c.Request().Header.VisitAll(func(key, value []byte) {
		headerName := string(key)
		if headerName != "Accept-Encoding" {
			req.Header.Set(string(key), string(value))
		}
	})

	req.Header.Set("X-Gateway-Key", "6F1ED002AB5595859014EBF0951522D9")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "ERROR_FORWARDING_REQUEST",
		})
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "ERROR_READ_RESPONSE",
		})
	}

	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(responseBody, &jsonResponse); err != nil {
		log.Println(jsonResponse)
		log.Printf("error Unmarshal: %v", err)
		log.Printf("Response Error: %s", resp.StatusCode)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "ERROR_JSON_RESPONSE",
		})
	}
	log.Printf("Response: %s", resp.StatusCode)
	return c.Status(int(resp.StatusCode)).JSON(jsonResponse)
}

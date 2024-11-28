package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

var backends = []struct {
	URL  string
	Code string
}{
	{"https://streambe01.indonesia.us.kg", "MAWAR"},
	{"https://partial-maureene-zaxkyu-237ae0c6.koyeb.app", "MELATI"},
	{
		"https://app-80a2ec7d-b99d-4abc-b771-24bf19f1b0a4.cleverapps.io",
		"CEMPAKA",
	},
}

type GatewayController struct{}

func getRandomBackend() (string, string) {
	rand.Seed(time.Now().UnixNano())
	backend := backends[rand.Intn(len(backends))]
	return backend.URL, backend.Code
}

func (gc *GatewayController) Index(c *fiber.Ctx) error {
	startTime := time.Now()
	method := c.Method()
	path := c.Path()
	body := c.Body()
	ip := c.Get("Cf-Connecting-Ip")

	log.Printf("Received Request :[ %s ] %s %s", ip, method, path)

	client := &http.Client{
		Timeout: 500 * time.Second,
	}

	backendUrl, backendCode := getRandomBackend()
	url := backendUrl + path

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
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

	// log.Println(req.Header)

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
	executionTime := time.Since(startTime).Round(time.Second / 10)
	jsonResponse["debug"] = fiber.Map{
		"response_by":    backendCode,
		"execution_time": executionTime.String(),
	}

	log.Printf(
		"Response: [%d | %s | %s | %s | %s]",
		resp.StatusCode,
		executionTime,
		backendCode,
		path,
		c.Get("User-Agent"),
	)

	return c.Status(int(resp.StatusCode)).JSON(jsonResponse)
}

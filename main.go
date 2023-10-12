package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"_message": "Hey! You've hit MicroDKIM. This is a service used by https://helium.ws to create DKIM keys.",
			"canIuse?": "Sure! Just send a GET request to https://dkim.helium.ws/new and you'll recieve a private and public key.",
			"error": false,
		})
	})

	app.Get("/new", func(c *fiber.Ctx) error {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"message": err.Error(),
			})
		}

		publicKey := &privateKey.PublicKey

		privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"message": err.Error(),
			})
		}

		publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"message": err.Error(),
			})
		}

		privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKeyDER)
		publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyDER)

		return c.JSON(fiber.Map{
			"error": false,
			"privateKey": privateKeyBase64,
			"publicKey": publicKeyBase64,
		})
	});

	app.Listen(":3000")
}
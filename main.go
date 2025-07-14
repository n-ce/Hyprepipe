package main

import (
	"fmt"
	"os"

	"codeberg.org/Hyperpipe/hyperpipe-backend/lib"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetHeaders(c *fiber.Ctx) error {
	c.Type("json", "utf-8")
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	c.Set("Access-Control-Max-Age", "300")

	return c.Next()
}

func HandleHealth(c *fiber.Ctx) error {
	return c.SendStatus(204)
}

func HandleExplore(c *fiber.Ctx) error {
	res, status := lib.GetExplore()

	return c.Status(status).JSON(res)
}

func HandleNext(c *fiber.Ctx) error {
	res, status := lib.GetNext(
		c.Params("id"),
		c.Query("queue") == "avoid")

	return c.Status(status).JSON(res)
}

func HandleGenres(c *fiber.Ctx) error {
	res, status := lib.GetGenres()

	return c.Status(status).JSON(res)
}

func HandleGenre(c *fiber.Ctx) error {
	res, status := lib.GetGenre(c.Params("id"))

	return c.Status(status).JSON(res)
}

func HandleCharts(c *fiber.Ctx) error {
	res, status := lib.GetCharts(
		c.Query("params"),
		c.Query("code"))

	return c.Status(status).JSON(res)
}

func HandleLyrics(c *fiber.Ctx) error {
	res, status := lib.GetLyrics(c.Params("id"))

	return c.Status(status).JSON(res)
}

func HandleAlbum(c *fiber.Ctx) error {
	res, status := lib.GetAlbumUrl(c.Params("id"))
	return c.Status(status).JSON(res)
}

func HandleArtist(c *fiber.Ctx) error {
	res, status := lib.GetArtist(c.Params("id"))

	return c.Status(status).JSON(res)
}

func HandleArtistNext(c *fiber.Ctx) error {
	if c.Params("id") == "" || c.Params("params") == "" ||
		c.Query("ct") == "" || c.Query("v") == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid params",
		})
	}

	res, status := lib.GetArtistNext(c.Params("id"),
		c.Params("params"),
		c.Query("ct"),
		c.Query("v"))

	return c.Status(status).JSON(res)
}

func main() {
	if os.Getenv("HYP_PROXY") == "" {
		fmt.Println("HYP_PROXY is empty!")
		os.Exit(1)
	}

	cfg := fiber.Config{
		Prefork: false,
	}

	if os.Getenv("HYP_PREFORK") == "1" {
		cfg.Prefork = true
	}

	app := fiber.New(cfg)

	app.Use(SetHeaders)
	app.Use(recover.New())

	app.Get("/healthz", HandleHealth)
	app.Get("/explore", HandleExplore)
	app.Get("/genres", HandleGenres)
	app.Get("/genres/:id", HandleGenre)
	app.Get("/charts", HandleCharts)
	app.Get("/next/:id", HandleNext)
	app.Get("/lyrics/:id", HandleLyrics)
	app.Get("/album/:id", HandleAlbum)
	app.Get("/channel/:id", HandleArtist)
	app.Get("/next/channel/:id/:params", HandleArtistNext)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Listen(":" + port)
}

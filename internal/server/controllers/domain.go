package controllers

import "github.com/gofiber/fiber/v3"

func (c *Controller) GetDomains(ctx fiber.Ctx) error {
	return ctx.JSON(c.config.Domains)
}

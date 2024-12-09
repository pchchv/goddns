package controllers

func (c *Controller) GetDomains(ctx *fiber.Ctx) error {
	return ctx.JSON(c.config.Domains)
}

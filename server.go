package hugo_auth

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/gofiber/fiber"
	"github.com/markbates/goth"

	"github.com/markbates/goth/providers/google"

	"github.com/gofiber/fiber/middleware"
	"github.com/shareed2k/goth_fiber"
)

func App() {
	godotenv.Load()
	app := fiber.New()

	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	goth.UseProviders(
		google.New(os.Getenv("OAUTH_KEY"), os.Getenv("OAUTH_SECRET"), os.Getenv("CALLBACK_URL")),
	)

	app.Use(func(c *fiber.Ctx) {
		googleAuth, _ := goth_fiber.GetFromSession("google", c)

		provider, _ := goth.GetProvider("google")
		sess, _ := provider.UnmarshalSession(googleAuth)
		user, _ := provider.FetchUser(sess)

		org, keyPresent := user.RawData["hd"]

		if keyPresent == true && org == os.Getenv("ORGANIZATION") {
			c.Next()
		} else {
			if c.Path() == "/login/google" || c.Path() == "/auth/callback" || c.Path() == "/logout" {
				c.Next()
			} else {
				c.Redirect("/login/google")
			}
		}
	})

	app.Static("/", "./public", fiber.Static{
		Browse: true,
	})

	app.Get("/login/:provider", goth_fiber.BeginAuthHandler)

	app.Get("/auth/callback", func(ctx *fiber.Ctx) {
		user, err := goth_fiber.CompleteUserAuth(ctx)
		if err != nil {
			ctx.Status(400).JSON(map[string]string{"message": err.Error()})
			return
		}

		if user.RawData["hd"] == os.Getenv("ORGANIZATION") {
			ctx.Redirect("/")
		} else {
			ctx.Redirect("/logout")
		}
	})

	app.Get("/logout", func(ctx *fiber.Ctx) {
		if err := goth_fiber.Logout(ctx); err != nil {
			log.Fatal(err)
		}

		ctx.Status(401).JSON(map[string]string{"message": "You are not authorized to see this content"})
	})

	if err := app.Listen(8088); err != nil {
		log.Fatal(err)
	}
}

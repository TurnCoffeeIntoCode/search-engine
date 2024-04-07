package routes

import (
	"coffeeintocode/search-engine/db"
	"coffeeintocode/search-engine/utils"
	"coffeeintocode/search-engine/views"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AdminClaims struct {
	User                 string `json:"user"`
	Id                   string `json:"id"`
	jwt.RegisteredClaims `json:"claims"`
}

func DashboardHandler(c *fiber.Ctx) error {
	settings := &db.SearchSettings{}
	err := settings.Get()
	if err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Something went wrong</h2>")
	}
	amount := strconv.FormatUint(uint64(settings.Amount), 10)
	return render(c, views.Home(amount, settings.SearchOn, settings.AddNew))
}

type settingsform struct {
	Amount   uint   `form:"amount"`
	SearchOn string `form:"searchOn"`
	AddNew   string `form:"addNew"`
}

func DashboardPostHandler(c *fiber.Ctx) error {
	input := settingsform{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Something went wrong</h2>")
	}
	// Convert checkbox 'on' values to boolean
	addNew := false
	if input.AddNew == "on" {
		addNew = true
	}
	searchOn := false
	if input.SearchOn == "on" {
		searchOn = true
	}
	settings := &db.SearchSettings{}
	settings.Amount = input.Amount
	settings.SearchOn = searchOn
	settings.AddNew = addNew
	err := settings.Update()
	if err != nil {
		fmt.Println(err)
		return c.SendString("<h2>Error: Something went wrong</h2>")
	}
	c.Append("HX-Refresh", "true")
	return c.SendStatus(200)
}

func LoginHandler(c *fiber.Ctx) error {
	return render(c, views.Login())
}

type loginform struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func LoginPostHandler(c *fiber.Ctx) error {
	input := loginform{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Something went wrong</h2>")
	}
	user := &db.User{}
	user, err := user.LoginAsAdmin(input.Email, input.Password)
	if err != nil {
		c.Status(401)
		c.Append("content-type", "text/html")
		return c.SendString("<h2>Error: Unauthorised</h2>")
	}

	signedToken, err := utils.CreateNewAuthToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		c.Status(500)
		return c.SendString("<h2>Error:Something went wrong logging in, please try again.</h2>")
	}

	// Create and set the cookie
	cookie := fiber.Cookie{
		Name:     "admin",
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true, // Meant only for the server
	}
	c.Cookie(&cookie)
	c.Append("HX-Redirect", "/")
	return c.SendStatus(200)
}

func LogoutHandler(c *fiber.Ctx) error {
	c.ClearCookie("admin")
	c.Set("HX-Redirect", "/login")
	return c.SendStatus(200)
}

func AuthMiddleware(c *fiber.Ctx) error {
	// Get the cookie by name
	cookie := c.Cookies("admin")
	if cookie == "" {
		return c.Redirect("/login", 302)
	}
	// Parse the cookie & check for errors
	token, err := jwt.ParseWithClaims(cookie, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return c.Redirect("/login", 302)
	}
	// Parse the custom claims & check jwt is valid
	_, ok := token.Claims.(*AdminClaims)
	if ok && token.Valid {
		return c.Next()
	}
	return c.Redirect("/login", 302)
}

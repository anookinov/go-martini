package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Add global methods if data is a map
	if 	viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

type User struct {
	Name string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

// e.GET("/users/:id", getUser)
func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

// e.GET("/show", show)
func show(c echo.Context) error {
	// Get team and member from the query string
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:" + team + ", member:" + member)
}

// e.POST("/save", save)
func save(c echo.Context) error {
	// Get name and email
	name := c.FormValue("name")
	email := c.FormValue("email")
	return c.String(http.StatusOK, "name:" + name + ", email:" + email)
}

// e.POST("/savefile", savefile)
func saveFile(c echo.Context) error {
	// Get name
	name := c.FormValue("name")
	//Get avatar
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return err
	}

	// Source
	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(avatar.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "<b>Thank you!" + name + "</b>")
}

func postUsers(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)
	// or
	// return c.XML(http.StatusCreated, u)
}

func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func helloTemplate(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}

func uriTemplate(c echo.Context) error {
	return c.Render(http.StatusOK, "template.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func basicAuth(username, password string, c echo.Context) (bool, error) {
	if username == "joe" && password == "secret" {
		return true, nil
	}
	return false, nil
}

func routeMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		println("request to /users")
		return next(c)
	}
}

func getUsers(c echo.Context) error {
	return c.String(http.StatusOK, "/users")
}

func main() {
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	
	e := echo.New()
	e.Renderer = renderer

	// e.Pre(middleware.HTTPSRedirect())
	// HTTPSRedirect
	// HTTPSWWWRedirect
	// WWWRedirect
	// NonWWWRedirect
	// AddTrailingSlash
	// RemoveTrailingSlash
	// MethodOverride
	// Rewrite

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// BodyLimit
	// Logger
	// Gzip
	// Recover
	// BasicAuth
	// JWTAuth
	// Secure
	// CORS
	// Static

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
		// 	if strings.HasPrefix(c.Request().Host, "localhost") {
		// 		return true
		// 	}
		// 	return false
			return strings.HasPrefix(c.Request().Host, "localhost")
		},
	}))

	// Group level middleware
	// g := e.Group("/admin", middleware.BasicAuth())
	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(basicAuth))

	// Route level middleware
	track := routeMiddleware

	// e.POST("/users", saveUser)
	e.GET("/users/:id", getUser)
	// e.PUT("/users/:id", updateUser)
	// e.DELETE("/users/:id", deleteUser)
	e.GET("/show", show)
	e.POST("/save", save)
	e.POST("/savefile", saveFile)
	e.POST("/users", postUsers)
	e.Static("/static", "assets")
	e.File("/", "public/index.html")
	e.File("/favicon.ico", "images/favicon.ico")
	e.File("/images/site.webmanifest", "images/site.webmanifest")
	e.GET("/hello", helloWorld)
	e.GET("/hellotemplate", helloTemplate)
	e.GET("/users", getUsers, track)

	// Named route "foobar"
	e.GET("/uritemplate", uriTemplate).Name = "foobar"
	e.Logger.Fatal(e.Start(":1323"))
}
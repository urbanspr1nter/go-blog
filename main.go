package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"log"
	"os"
	"strings"
)

type Post struct {
	Id        int32  `json:"id"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Created   string `json:"created"`
	Published bool   `json:"published"`
}

type Configuration struct {
	Title          string `json:"title"`
	Author         string `json:"author"`
	DatabaseUrlEnv string `json:"database_url_env"`
	GithubUrl      string `json:"github_url"`
	LinkedInUrl    string `json:"linkedin_url"`
	HttpPort       string `json:"http_port"`
}

/*
Utilizes the existing MySQL connection to query for all published posts ordered by date created.
*/
func getPosts(conn *sql.DB) []Post {
	listSlice := make([]Post, 0, 10)

	rows, err := conn.Query("select * from post where published = 1 order by created desc")
	if err != nil {
		log.Fatal("Could not make a connection to the MySQL database.")
	}

	for rows.Next() {
		var id int32
		var slug string
		var title string
		var content string
		var created string
		var published bool

		err := rows.Scan(&id, &slug, &title, &content, &created, &published)
		if err != nil {
			log.Fatal("One of the columns in the row could not be mapped and assigned to a variable.")
		}

		post := Post{id, slug, title, content, created, published}

		listSlice = append(listSlice, post)
	}

	return listSlice
}

/*
Gets the configuration.json file at the same directory of the app and deserializes it to a configuration struct
*/
func getConfiguration() Configuration {
	data, err := os.ReadFile(makeFullPath("configuration.json"))
	if err != nil {
		log.Fatal("Could not read configuration.json. Does it exist in the same directory as the app?")
	}

	var configuration Configuration
	err = json.Unmarshal(data, &configuration)
	if err != nil {
		log.Fatal("Could not deserialize JSON to a Configuration structure")
	}

	return configuration
}

/* Handy helper to always build render template variables so that pages ender consistently with the correct data */
func buildTemplateArguments(params fiber.Map, configuration *Configuration) fiber.Map {
	result := fiber.Map{
		"PageTitle":   configuration.Title,
		"GitHubUrl":   configuration.GithubUrl,
		"LinkedInUrl": configuration.LinkedInUrl,
		"Author":      configuration.Author,
	}

	for k, v := range params {
		result[k] = v
	}

	return result
}

func makeFullPath(relativePath string) string {
	pwd, err := os.Getwd()

	if err != nil {
		log.Fatal("Cannot get the current working directory.")
	}

	return fmt.Sprintf("%s/%s", pwd, relativePath)
}

func main() {
	configuration := getConfiguration()

	// `database_url_env` value is actually the environment variable name.
	//
	// If the database_url_env value in configuration.json is "DATABASE_URL"
	// then os.GetEnv("DATABASE_URL") will contain the real endpoint.
	databaseUrlEnvVar := os.Getenv(configuration.DatabaseUrlEnv)
	conn, err := sql.Open("mysql", databaseUrlEnvVar)
	if err != nil {
		log.Fatal("Unable to connect to the database. Either this machine is offline, or the server is offline.")
	}

	defer conn.Close()

	// Initialize the view engine and set the directory to /views
	viewEngine := html.New(makeFullPath("views"), ".html")
	app := fiber.New(fiber.Config{
		Views: viewEngine,
	})

	// Register static path to be at the root. The files will be found in /public.
	// So if a request comes in for /picture.jpg, then the server will attempt to find
	// picture.jpg within the /public folder locally. Ex: /public/picture.jpg.
	app.Static("/", makeFullPath("public"))

	app.Get("/", func(c *fiber.Ctx) error {
		listSlice := getPosts(conn)
		return c.Render("index", buildTemplateArguments(fiber.Map{
			"Posts": listSlice,
		}, &configuration), "layout")
	})

	app.Get("/about", func(c *fiber.Ctx) error {
		return c.Render("about", buildTemplateArguments(fiber.Map{}, &configuration), "layout")
	})

	app.Get("/resume", func(c *fiber.Ctx) error {
		return c.Render("resume", buildTemplateArguments(fiber.Map{}, &configuration), "layout")
	})

	app.Get("/p/:slug", func(c *fiber.Ctx) error {
		userSlug := strings.ToLower(c.Params("slug"))

		// Really lame way to sanitize for now
		invalidCharacters := []string{
			"!", "#", "$", "%", "^",
			"&", "*", "(", ")", " ",
			";", ":", "\"", "\\", "/",
			"'", "?", ".", "<", ">",
			"{", "}", "[", "]",
		}

		for _, invalidCharacter := range invalidCharacters {
			userSlug = strings.ReplaceAll(userSlug, invalidCharacter, "")
		}

		query := fmt.Sprintf("select * from post where slug = '%s' and published = 1 order by created desc", userSlug)
		rows, err := conn.Query(query)
		if err != nil {
			return c.SendStatus(404)
		}

		if rows.Next() {
			var id int32
			var slug string
			var title string
			var content string
			var created string
			var published bool

			err := rows.Scan(&id, &slug, &title, &content, &created, &published)
			if err != nil {
				c.SendStatus(500)
			}

			post := Post{id, slug, title, content, created, published}

			return c.Render("post", buildTemplateArguments(fiber.Map{
				"Title":   post.Title,
				"Content": post.Content,
				"Created": post.Created,
			}, &configuration), "layout")
		}

		// Otherwise we couldn't find anything.
		return c.SendStatus(404)
	})

	listenErr := app.Listen(fmt.Sprintf("0.0.0.0:%s", configuration.HttpPort))
	if listenErr != nil {
		log.Fatal("Could not bind to the port.")
	}
}

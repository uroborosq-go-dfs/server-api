package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/uroborosq-go-dfs/server/server"
	"log"
	"os"
)

// @title Go-DFS
// @version 1.0
// @description Go-DFS server api
// @host localhost:3000
// @BasePath /
func main() {
	app := fiber.New()

	port := os.Getenv("GODFS_SERVER_API_PORT")
	if port == "" {
		port = ":3000"
	}

	host := os.Getenv("GODFS_DB_HOST")
	dbPort := os.Getenv("GODFS_DB_PORT")
	user := os.Getenv("GODFS_DB_USER")
	password := os.Getenv("GODFS_DB_PASSWORD")
	dbname := os.Getenv("GODFS_DB_DBNAME")
	driver := os.Getenv("GODFS_DB_DRIVER")

	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, dbPort, user, password, dbname)
	s, err := server.CreateServer(driver, conn)
	if err != nil {
		log.Fatal(err.Error())
	}

	c := New(s)

	app.Post("/node/add", c.HandleAddNode)
	app.Delete("/node/remove", c.HandleDeleteNode)
	app.Patch("/node/clean", c.HandleCleanNode)
	app.Get("/list", c.HandleListAllFiles)
	app.Get("/list/node", c.HandleListOfNode)
	app.Get("/file/get", c.HandleGetFile)
	app.Post("/file/add", c.HandleAddFile)
	app.Delete("/file/delete", c.HandleRemoveFile)

	log.Fatal(app.Listen(port))
}

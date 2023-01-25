package main

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	guid "github.com/google/uuid"
	"github.com/uroborosq-go-dfs/server/connector"
	"github.com/uroborosq-go-dfs/server/server"
	"log"
	"net/http"
	"os"
	"strconv"
)

func New(s *server.Server) *Controller {
	return &Controller{s: s}
}

type Controller struct {
	s *server.Server
}

func (s *Controller) HandleAddNode(c *fiber.Ctx) error {
	url := c.Query("url")
	maxSize, err := strconv.ParseInt(c.Query("maximum_size"), 0, 64)
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		return err
	}
	connectionType, err := strconv.Atoi(c.Query("connection_type"))
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		return err
	}

	id, err := s.s.AddNode(url, "nope", maxSize, connector.NetConnectorType(connectionType))
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	c.SendStatus(http.StatusOK)
	return c.SendString(id.String())
}

func (s *Controller) HandleDeleteNode(c *fiber.Ctx) error {
	id, err := guid.Parse(c.Query("id"))
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		return errors.New("query argument path is empty")
	}

	err = s.s.RemoveNode(id)
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (s *Controller) HandleCleanNode(c *fiber.Ctx) error {
	id, err := guid.Parse(c.Query("id"))
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		return errors.New("query argument path is empty")
	}

	err = s.s.CleanNode(id)
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (s *Controller) HandleListAllFiles(c *fiber.Ctx) error {
	paths, sizes, err := s.s.ListOfAllFiles()
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		return err
	}
	filePaths := make([]*FilePath, len(paths))
	for i := 0; i < len(paths); i++ {
		filePaths[i] = &FilePath{
			Path: paths[i],
			Size: sizes[i],
		}
	}

	jsonStr, err := json.Marshal(filePaths)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	_ = c.Send(jsonStr)

	return c.SendStatus(http.StatusOK)
}

func (s *Controller) HandleListOfNode(c *fiber.Ctx) error {
	id, err := guid.Parse(c.Query("id"))
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		return err
	}
	paths, sizes, err := s.s.ListOfNodeFiles(id)
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		return err
	}
	filePaths := make([]*FilePath, len(paths))
	for i := 0; i < len(paths); i++ {
		filePaths[i] = &FilePath{
			Path: paths[i],
			Size: sizes[i],
		}
	}

	jsonStr, err := json.Marshal(filePaths)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	_ = c.Send(jsonStr)

	return c.SendStatus(http.StatusOK)
}

func (s *Controller) HandleGetFile(c *fiber.Ctx) error {
	filePath := c.Query("path")
	if filePath == "" {
		c.SendStatus(http.StatusBadRequest)
		return errors.New("query argument is empty")
	}

	err := s.s.GetFile(filePath, filePath)

	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}

	err = c.SendFile(filePath)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}

	err = os.Remove(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	return c.SendStatus(http.StatusOK)
}

func (s *Controller) HandleAddFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		return err
	}
	err = c.SaveFile(file, file.Filename)
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		return err
	}
	c.SendStatus(http.StatusOK)
	return s.s.AddFile(file.Filename, file.Filename)
}

func (s *Controller) HandleRemoveFile(c *fiber.Ctx) error {
	path := c.Query("path")
	if path == "" {
		c.SendStatus(http.StatusBadRequest)
		return errors.New("query argument path is empty")
	}
	err := s.s.RemoveFile(path)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}

	return nil
}

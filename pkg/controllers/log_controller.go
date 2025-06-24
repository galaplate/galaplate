package controllers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sheenazien8/goplate/env"
	"github.com/sheenazien8/goplate/logs"
)

type LogFile struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	ModifiedTime time.Time `json:"modified_time"`
	Content      string    `json:"content,omitempty"`
}

type LogController struct {
}

func (c *LogController) ShowLogsPage(ctx *fiber.Ctx) error {
	var logsDir = "./storage/logs"
	var selectedFile = ctx.Query("file")
	var selectedContent = ""

	files, err := ioutil.ReadDir(logsDir)
	if err != nil {
		logs.Error("Unable to read logs directory:", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Unable to read logs directory")
	}

	var logFiles []LogFile
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".log") {
			logFiles = append(logFiles, LogFile{
				Name:         file.Name(),
				Size:         file.Size(),
				ModifiedTime: file.ModTime(),
			})
		}
	}

	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].ModifiedTime.After(logFiles[j].ModifiedTime)
	})

	if selectedFile != "" {
		if !strings.HasSuffix(selectedFile, ".log") {
			return ctx.Status(fiber.StatusBadRequest).SendString("Invalid file type")
		}

		var logPath = filepath.Join(logsDir, selectedFile)

		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			return ctx.Status(fiber.StatusNotFound).SendString("Log file not found")
		}

		content, err := ioutil.ReadFile(logPath)
		if err != nil {
			logs.Error("Unable to read log file:", err)
			return ctx.Status(fiber.StatusInternalServerError).SendString("Unable to read log file")
		}

		selectedContent = string(content)
	}

    appName := env.Get("APP_NAME")

	return ctx.Render("logs", fiber.Map{
		"Title":           fmt.Sprintf("%s - Logs Viewer", appName),
		"LogFiles":        logFiles,
		"SelectedFile":    selectedFile,
		"SelectedContent": selectedContent,
	})
}

var LogControllerInstance = &LogController{}

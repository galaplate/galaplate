package controllers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/galaplate/core/env"
	"github.com/galaplate/core/logger"
	"github.com/gofiber/fiber/v2"
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
		logger.Error("LogController@ShowLogsPage", map[string]any{
			"messages": fmt.Sprintf("Unable to read logs directory: %s", err.Error()),
		})
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
			logger.Error("LogController@ShowLogsPage",map[string]any{
                "messages": fmt.Sprintf("Unable to read log file: %s", err.Error()),
            })
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

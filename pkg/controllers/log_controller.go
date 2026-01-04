package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LogController struct{}

type LogEntry struct {
	Timestamp      string         `json:"timestamp"`
	Level          string         `json:"level"`
	Message        string         `json:"message"`
	AdditionalInfo map[string]any `json:"additional_info,omitempty"`
}

type LogViewerData struct {
	LogFiles    []string
	CurrentFile string
	Logs        []LogEntryWithJSON
	LogsJSON    string
	TotalLogs   int
	ErrorCount  int
	InfoCount   int
	CurrentPage int
	TotalPages  int
	PageSize    int
	HasPrevious bool
	HasNext     bool
}

type LogEntryWithJSON struct {
	Timestamp          string         `json:"timestamp"`
	Level              string         `json:"level"`
	Message            string         `json:"message"`
	AdditionalInfo     map[string]any `json:"additional_info,omitempty"`
	AdditionalInfoJSON string         `json:"-"`
}

func (lvc *LogController) Export(c *fiber.Ctx) error {
	logDir := "./storage/logs"

	logFiles, err := getLogFiles(logDir)
	if err != nil {
		return c.Status(500).SendString("Error reading log directory")
	}

	currentFile := c.Query("file")
	if currentFile == "" || len(logFiles) == 0 {
		return c.Status(400).SendString("No log file specified")
	}

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	format := c.Query("format")

	logs, err := parseLogFile(filepath.Join(logDir, currentFile))
	if err != nil {
		return c.Status(500).SendString("Error reading log file")
	}

	filteredLogs := filterLogsByDate(logs, dateFrom, dateTo)

	if format == "json" {
		c.Set("Content-Type", "application/json")
		c.Set("Content-Disposition", `attachment; filename="logs-export.json"`)
		return c.JSON(fiber.Map{
			"exported_at": time.Now().Format(time.RFC3339),
			"file":        currentFile,
			"total_logs":  len(filteredLogs),
			"logs":        filteredLogs,
		})
	}

	if format == "csv" {
		csvContent := "Timestamp,Level,Message,ActionField,AdditionalInfo\n"
		for _, log := range filteredLogs {
			additionalInfoStr := ""
			if log.AdditionalInfo != nil {
				data, _ := json.Marshal(log.AdditionalInfo)
				additionalInfoStr = string(data)
			}

			action := ""
			if log.AdditionalInfo != nil {
				if a, ok := log.AdditionalInfo["action"]; ok {
					action = fmt.Sprintf("%v", a)
				}
			}

			csvLine := fmt.Sprintf("%q,%q,%q,%q,%q\n",
				log.Timestamp,
				log.Level,
				log.Message,
				action,
				additionalInfoStr,
			)
			csvContent += csvLine
		}

		c.Set("Content-Type", "text/csv")
		c.Set("Content-Disposition", `attachment; filename="logs-export.csv"`)
		return c.SendString(csvContent)
	}

	return c.Status(400).SendString("Invalid export format. Use 'json' or 'csv'")
}

func (lvc *LogController) CleanupLogs(c *fiber.Ctx) error {
	logDir := "./storage/logs"
	daysStr := c.Query("days", "30")
	days := 30
	if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
		days = d
	}

	cutoffDate := time.Now().AddDate(0, 0, -days)
	var deletedCount int
	var totalSize int64

	err := filepath.WalkDir(logDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if info.ModTime().Before(cutoffDate) {
			totalSize += info.Size()
			if err := os.Remove(path); err == nil {
				deletedCount++
			}
		}

		return nil
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to cleanup logs",
		})
	}

	return c.JSON(fiber.Map{
		"success":        true,
		"deleted_count":  deletedCount,
		"total_size_mb":  float64(totalSize) / (1024 * 1024),
		"cutoff_date":    cutoffDate.Format("2006-01-02"),
		"retention_days": days,
	})
}

func (lvc *LogController) GetLogStats(c *fiber.Ctx) error {
	logDir := "./storage/logs"

	var totalFiles int
	var totalSize int64
	var oldestFile time.Time
	var newestFile time.Time

	err := filepath.WalkDir(logDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		totalFiles++
		info, err := d.Info()
		if err != nil {
			return nil
		}

		totalSize += info.Size()
		modTime := info.ModTime()

		if oldestFile.IsZero() || modTime.Before(oldestFile) {
			oldestFile = modTime
		}
		if newestFile.IsZero() || modTime.After(newestFile) {
			newestFile = modTime
		}

		return nil
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to get log statistics",
		})
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"total_files":   totalFiles,
		"total_size_mb": float64(totalSize) / (1024 * 1024),
		"oldest_date":   oldestFile.Format("2006-01-02"),
		"newest_date":   newestFile.Format("2006-01-02"),
	})
}

func (lvc *LogController) Index(c *fiber.Ctx) error {
	logDir := "./storage/logs"

	logFiles, err := getLogFiles(logDir)
	if err != nil {
		return c.Status(500).SendString("Error reading log directory")
	}

	if len(logFiles) == 0 {
		return c.Status(404).SendString("No log files found")
	}

	currentFile := c.Query("file")
	if currentFile == "" {
		currentFile = logFiles[0]
	}

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	pageSize := 50
	pageSizeParam := c.Query("page_size")
	if pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 500 {
			pageSize = ps
		}
	}

	page := 1
	pageParam := c.Query("page")
	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	logs, err := parseLogFile(filepath.Join(logDir, currentFile))
	if err != nil {
		return c.Status(500).SendString("Error reading log file")
	}

	filteredLogs := filterLogsByDate(logs, dateFrom, dateTo)

	totalLogs := len(filteredLogs)
	errorCount := 0
	infoCount := 0

	for _, log := range filteredLogs {
		switch log.Level {
		case "error":
			errorCount++
		case "info":
			infoCount++
		}
	}

	totalPages := (totalLogs + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	if page > totalPages {
		page = totalPages
	}

	startIdx := (page - 1) * pageSize
	endIdx := startIdx + pageSize
	if endIdx > totalLogs {
		endIdx = totalLogs
	}

	paginatedLogs := filteredLogs[startIdx:endIdx]

	logsWithJSON := make([]LogEntryWithJSON, len(paginatedLogs))
	for i, log := range paginatedLogs {
		jsonBytes, _ := json.Marshal(log.AdditionalInfo)
		logsWithJSON[i] = LogEntryWithJSON{
			Timestamp:          log.Timestamp,
			Level:              log.Level,
			Message:            log.Message,
			AdditionalInfo:     log.AdditionalInfo,
			AdditionalInfoJSON: string(jsonBytes),
		}
	}

	logsJSON, _ := json.Marshal(paginatedLogs)

	data := LogViewerData{
		LogFiles:    logFiles,
		CurrentFile: currentFile,
		Logs:        logsWithJSON,
		LogsJSON:    string(logsJSON),
		TotalLogs:   totalLogs,
		ErrorCount:  errorCount,
		InfoCount:   infoCount,
		CurrentPage: page,
		TotalPages:  totalPages,
		PageSize:    pageSize,
		HasPrevious: page > 1,
		HasNext:     page < totalPages,
	}

	tmpl, err := template.ParseFiles("templates/log-viewer.html")
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	c.Type("html")
	return tmpl.Execute(c.Response().BodyWriter(), data)
}

func getLogFiles(logDir string) ([]string, error) {
	var logFiles []string

	err := filepath.WalkDir(logDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".log") {
			logFiles = append(logFiles, d.Name())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(logFiles)))

	return logFiles, nil
}

func parseLogFile(filePath string) ([]LogEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var logs []LogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var logEntry LogEntry
		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			continue
		}

		if logEntry.Timestamp != "" {
			if parsedTime, err := time.Parse(time.RFC3339, logEntry.Timestamp); err == nil {
				logEntry.Timestamp = parsedTime.Format("2006-01-02 15:04:05")
			}
		}

		if logEntry.Level == "" {
			logEntry.Level = "info"
		} else {
			logEntry.Level = strings.ToLower(logEntry.Level)
		}

		logs = append(logs, logEntry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	reverseSlice(logs)

	return logs, nil
}

func reverseSlice(logs []LogEntry) {
	for i := 0; i < len(logs)/2; i++ {
		j := len(logs) - 1 - i
		logs[i], logs[j] = logs[j], logs[i]
	}
}

func filterLogsByDate(logs []LogEntry, dateFrom, dateTo string) []LogEntry {
	if dateFrom == "" && dateTo == "" {
		return logs
	}

	var filteredLogs []LogEntry

	for _, log := range logs {
		if log.Timestamp == "" {
			continue
		}

		logTime, err := time.Parse("2006-01-02 15:04:05", log.Timestamp)
		if err != nil {
			continue
		}

		include := true

		if dateFrom != "" {
			fromTime, err := time.Parse("2006-01-02", dateFrom)
			if err == nil && logTime.Before(fromTime) {
				include = false
			}
		}

		if dateTo != "" && include {
			toTime, err := time.Parse("2006-01-02", dateTo)
			if err == nil {
				toTime = toTime.Add(24*time.Hour - time.Second)
				if logTime.After(toTime) {
					include = false
				}
			}
		}

		if include {
			filteredLogs = append(filteredLogs, log)
		}
	}

	return filteredLogs
}

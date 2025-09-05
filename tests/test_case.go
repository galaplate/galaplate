package tests

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/galaplate/core/bootstrap"
	"github.com/galaplate/galaplate/router"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// RefreshDatabase interface for tests that need database refresh
type RefreshDatabase interface {
	RefreshDatabase() error
}

type TestCase struct {
	suite.Suite
	App               *fiber.App
	DB                *gorm.DB
	refreshDatabase   bool
	databaseRefreshed bool
}

func (tc *TestCase) SetupTest() {
	// Always ensure we're in project root for the entire test setup
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..")
	absProjectRoot, _ := filepath.Abs(projectRoot)

	// Change to project root and stay there for the test
	os.Chdir(absProjectRoot)

	if err := godotenv.Load(".env.testing"); err != nil {
		log.Panicf("Warning: Error loading .env file: %v", err)
	}

	// Set environment for testing the default is using .env file
	os.Setenv("APP_ENV", "testing")
	// os.Setenv("DB_DATABASE", "galaplate")

	// Refresh database if enabled and not already refreshed
	if tc.refreshDatabase && !tc.databaseRefreshed {
		if err := tc.RefreshDatabase(); err != nil {
			log.Printf("Warning: Failed to refresh database: %v", err)
			log.Printf("Continuing with existing database state...")
		} else {
			tc.databaseRefreshed = true
		}
	}

	// Use bootstrap to create a properly configured app for testing
	cfg := bootstrap.DefaultConfig()
	cfg.SetupRoutes = router.SetupRouter
	app := bootstrap.App(cfg)

	tc.App = app

	// Note: We intentionally don't restore the working directory
	// This ensures storage and other paths are relative to project root
}

// EnableRefreshDatabase enables database refresh for this test case
func (tc *TestCase) EnableRefreshDatabase() {
	tc.refreshDatabase = true
}

// RefreshDatabase runs the db:fresh command to refresh the test database
func (tc *TestCase) RefreshDatabase() error {
	cmd := exec.Command("go", "run", "main.go", "console", "db:fresh", "--force")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// RefreshDatabaseBetweenTests refreshes database before each test method
func (tc *TestCase) RefreshDatabaseBetweenTests() {
	// Reset the flag so database gets refreshed for each test
	tc.databaseRefreshed = false
}

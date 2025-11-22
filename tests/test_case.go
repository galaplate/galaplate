package tests

import (
	"io"
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

type TestCase struct {
	suite.Suite
	App               *fiber.App
	DB                *gorm.DB
	isRefreshDatabase bool
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

	// Refresh database if enabled and not already refreshed
	if tc.isRefreshDatabase && !tc.databaseRefreshed {
		if err := tc.refreshDatabase(); err != nil {
			log.Printf("Warning: Failed to refresh database: %v", err)
			log.Printf("Continuing with existing database state...")
		} else {
			tc.databaseRefreshed = true
		}
	}

	// Use bootstrap to create a properly configured app for testing
	// cfg := bootstrap.DefaultConfig()
	// cfg.SetupRoutes = router.SetupRouter
	app := bootstrap.NewApp(func(ac *bootstrap.AppConfig) {
		ac.SetupRoutes = router.SetupRouter
	})

	tc.App = app

	// Note: We intentionally don't restore the working directory
	// This ensures storage and other paths are relative to project root
}

// enableRefreshDatabase enables database refresh for this test case
func (tc *TestCase) enableRefreshDatabase() {
	tc.isRefreshDatabase = true
}

// refreshDatabase runs the db:fresh command to refresh the test database
func (tc *TestCase) refreshDatabase() error {
	cmd := exec.Command("go", "run", "main.go", "console", "db:fresh", "--force")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	return cmd.Run()
}

// seedDatabase runs the database seeders
func (tc *TestCase) seedDatabase() error {
	cmd := exec.Command("go", "run", "main.go", "console", "db:seed")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// refreshDatabaseBetweenTests refreshes database before each test method
func (tc *TestCase) refreshDatabaseBetweenTests() {
	tc.databaseRefreshed = false
}

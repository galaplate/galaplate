package tests

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/galaplate/core/bootstrap"
	"github.com/galaplate/galaplate/router"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestCase struct {
	suite.Suite
	App *fiber.App
	DB  *gorm.DB
}

func (tc *TestCase) SetupTest() {
	// Always ensure we're in project root for the entire test setup
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..")
	absProjectRoot, _ := filepath.Abs(projectRoot)

	// Change to project root and stay there for the test
	os.Chdir(absProjectRoot)

	// Set environment for testing the default is using .env file
	os.Setenv("APP_ENV", "testing")
	// os.Setenv("DB_DATABASE", "galaplate")

	// Use bootstrap to create a properly configured app for testing
	cfg := bootstrap.DefaultConfig()
	cfg.SetupRoutes = router.SetupRouter
	app := bootstrap.App(cfg)

	tc.App = app

	// Note: We intentionally don't restore the working directory
	// This ensures storage and other paths are relative to project root
}

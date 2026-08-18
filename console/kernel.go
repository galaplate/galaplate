package console

import (
	"github.com/galaplate/core/console"
	"github.com/galaplate/core/console/commands"
	"github.com/galaplate/galaplate/router"
	// "github.com/galaplate/galaplate/console/commands"
)

func RegisterCommands(kernel *console.Kernel) {
	// Register your custom console commands here
	// Example:
	// kernel.Register(&commands.SendwelcomeemailcommandCommand{})
	// Inject app-specific router into the built-in route:list command
	if cmd, ok := kernel.GetCommands()["route:list"]; ok {
		if routeCmd, ok := cmd.(*commands.RouteListCommand); ok {
			routeCmd.SetupRoutes = router.SetupRouter
		}
	}
}

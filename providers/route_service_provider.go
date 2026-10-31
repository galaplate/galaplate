package providers

import (
	"github.com/galaplate/core/fiber"
	"github.com/galaplate/core/provider"
	"github.com/galaplate/galaplate/router"
)

type RouteServiceProvider struct{}

func init() { provider.Add(&RouteServiceProvider{}) }

func (p *RouteServiceProvider) Register() {}

func (p *RouteServiceProvider) Boot() {}

func (p *RouteServiceProvider) Routes(app *fiber.App) {
	router.SetupRouter(app)
}

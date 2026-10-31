package providers

import (
	"github.com/galaplate/core/provider"
	"github.com/galaplate/galaplate/pkg/middleware"
)

type AppServiceProvider struct{}

func init() { provider.Add(&AppServiceProvider{}) }

func (p *AppServiceProvider) Register() {
	provider.AliasMiddleware("auth", middleware.JWTAuth())
}

func (p *AppServiceProvider) Boot() {}

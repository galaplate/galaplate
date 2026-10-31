package providers

import "github.com/galaplate/core/provider"

type ConsoleServiceProvider struct{}

func init() { provider.Add(&ConsoleServiceProvider{}) }

func (p *ConsoleServiceProvider) Register() {}

func (p *ConsoleServiceProvider) Boot() {}

func (p *ConsoleServiceProvider) Commands() []provider.Command {
	return []provider.Command{}
}

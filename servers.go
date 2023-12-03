package openengine

import "github.com/tahersoft-go/openengine/engine"

func (p *openEngine) AddServers(servers engine.ApiServers) OpenEngine {
	p.Servers = servers
	return p
}

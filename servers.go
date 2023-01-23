package openengine

import "gitlab.hoitek.fi/healthcare/services/maja/openengine/engine"

func (p *openEngine) AddServers(servers engine.ApiServers) OpenEngine {
	p.Servers = servers
	return p
}

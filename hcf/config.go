package hcf

import "github.com/df-mc/dragonfly/server"

// Config ...
type Config struct {
	Factions struct {
		Map   int
		Men   int
		Level int
	}
	server.UserConfig
}

// DefaultConfig ...
func DefaultConfig() Config {
	c := Config{
		UserConfig: server.DefaultConfig(),
	}
	c.Factions.Map = 1
	c.Factions.Men = 15
	c.Factions.Level = 2
	return c
}

package main

import (
	config "github.com/Fepozopo/gator/internal/config"
)

// state hold application-level state, such as the config.
type state struct {
	config *config.Config
}

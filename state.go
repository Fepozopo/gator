package main

import (
	config "github.com/Fepozopo/gator/internal/config"
	database "github.com/Fepozopo/gator/internal/database"
)

// state holds application-level state.
type state struct {
	db  *database.Queries
	cfg *config.Config
}

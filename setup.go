package main

import (
	"database/sql"
	"fmt"

	"github.com/ayrtonvitor/gator/internal/config"
	"github.com/ayrtonvitor/gator/internal/database"
	"github.com/ayrtonvitor/gator/internal/state"
)

func setup() (*state.State, error) {
	conf, err := config.Read()
	if err != nil {
		return nil, fmt.Errorf("Could not read %v", err)
	}
	state := &state.State{
		Config: &conf,
	}

	queries, err := getQueries(state.Config.ConnString)
	if err != nil {
		return nil, fmt.Errorf("Could not setup state: %w", err)
	}
	state.Db = queries

	return state, nil
}

func getQueries(connString string) (*database.Queries, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to DB to get queries: %w", err)
	}
	return database.New(db), nil
}

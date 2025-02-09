package state

import (
	"github.com/ayrtonvitor/gator/internal/config"
	"github.com/ayrtonvitor/gator/internal/database"
)

type State struct {
	Config *config.Config
	Db     *database.Queries
}

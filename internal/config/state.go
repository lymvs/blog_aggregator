package config

import (
	"github.com/lymvs/blog_aggregator/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *Config
}

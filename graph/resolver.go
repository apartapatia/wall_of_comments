package graph

import (
	"github.com/apartapatia/wall_of_comments/internal/database"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repo database.Repo
}

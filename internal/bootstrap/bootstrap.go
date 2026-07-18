package bootstrap

import (
	"database/sql"
	"errors"
)

var ErrDatabaseRequired = errors.New("database is required")

type Dependencies struct {
	Database *sql.DB
}

type Runtime struct {
	Database *sql.DB
}

func New(dependencies Dependencies) (*Runtime, error) {
	if dependencies.Database == nil {
		return nil, ErrDatabaseRequired
	}

	return &Runtime{Database: dependencies.Database}, nil
}

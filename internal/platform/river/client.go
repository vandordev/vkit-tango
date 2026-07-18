package river

import (
	"database/sql"

	riverqueue "github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
)

func NewClient(database *sql.DB) (*riverqueue.Client[*sql.Tx], error) {
	return riverqueue.NewClient(riverdatabasesql.New(database), &riverqueue.Config{
		Workers: riverqueue.NewWorkers(),
	})
}

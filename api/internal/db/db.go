package db

import (
	"database/sql"
)

type Driver struct {
	*sql.DB
}

package data

import (
	"github.com/jmoiron/sqlx"
)

type Passive struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type PassiveModel struct {
	*sqlx.DB
}

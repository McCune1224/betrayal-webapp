package data

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Ability struct {
	ID           int            `db:"id"`
	Name         string         `db:"name"`
	Description  string         `db:"description"`
	Charges      int            `db:"charges"`
	AnyAbility   bool           `db:"any_ability"`
	RoleSpecific string         `db:"role_specific"`
	Categories   pq.StringArray `db:"categories"`
}

// psql statement to add categories to ability
// ALTER TABLE abilities ADD COLUMN categories text[] DEFAULT '{}';

type AbilityModel struct {
	*sqlx.DB
}

package data

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Alliance struct {
	ID          int            `db:"id"`
	Name        string         `db:"name"`
	GameID      string         `db:"game_id"`
	Description string         `db:"description"`
	Members     pq.StringArray `db:"members"`
	Color       string         `db:"color"`
}

// psql to make name unique:
// ALTER TABLE alliances ADD CONSTRAINT unique_name UNIQUE (name);

type AllianceModel struct {
	*sqlx.DB
}

func (m *AllianceModel) Get(id int) (*Alliance, error) {
	var alliance *Alliance
	err := m.DB.Get(&alliance, "SELECT * FROM alliances WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return alliance, nil
}

func (m *AllianceModel) GetByName(name string) (*Alliance, error) {
	var alliance *Alliance
	err := m.DB.Get(&alliance, "SELECT * FROM alliances WHERE name=$1", name)
	if err != nil {
		return nil, err
	}
	return alliance, nil
}

func (m *AllianceModel) GetByMember(member string) ([]*Alliance, error) {
	var alliances []*Alliance
	err := m.Select(&alliances, "SELECT * FROM alliances WHERE $1 = ANY(members)", member)
	if err != nil {
		return nil, err
	}
	return alliances, nil
}

func (m *AllianceModel) Delete(id int) error {
	_, err := m.Exec("DELETE FROM alliances WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func (m *AllianceModel) Create(a *Alliance) error {
	_, err := m.Exec("INSERT INTO alliances (name, description, members, color, game_id) VALUES ($1, $2, $3, $4, $5)", a.Name, a.Description, a.Members, a.Color, a.GameID)
	if err != nil {
		return err
	}
	return nil
}

func (m *AllianceModel) CreateWithNameAndGameID(name, gameID string) error {
	_, err := m.Exec("INSERT INTO alliances (name, game_id) VALUES ($1, $2)", name, gameID)
	if err != nil {
		return err
	}
	return nil
}

func (m *AllianceModel) Update(a *Alliance) error {
	_, err := m.Exec("UPDATE alliances SET name=$1, description=$2, members=$3, color=$4, game_id=$5 WHERE id=$6", a.Name, a.Description, a.Members, a.Color, a.GameID, a.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *AllianceModel) UpdateMembers(id int, members pq.StringArray) error {
	_, err := m.Exec("UPDATE alliances SET members=$1 WHERE id=$2", members, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *AllianceModel) GetAll() ([]*Alliance, error) {
	var alliances []*Alliance
	err := m.Select(&alliances, "SELECT * FROM alliances")
	if err != nil {
		return nil, err
	}
	return alliances, nil
}

func (m *AllianceModel) GetAllByGame(gameID string) ([]*Alliance, error) {
	var alliances []*Alliance
	err := m.Select(&alliances, "SELECT * FROM alliances WHERE game_id=$1", gameID)
	if err != nil {
		return nil, err
	}
	return alliances, nil
}

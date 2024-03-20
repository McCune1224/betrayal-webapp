package data

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Player struct {
	ID                int    `db:"id"`
	Name              string `db:"name"`
	GameID            string `db:"game_id"`
	RoleID            int    `db:"role_id"`
	Alive             bool   `db:"alive"`
	Seat              int    `db:"seat"`
	Luck              int    `db:"luck"`
	LuckModifier      int    `db:"luck_modifier"`
	LuckStatus        string `db:"luck_status"`
	AlignmentOverride string `db:"alignment_override"`
	CreatedAt         string `db:"created_at"`
	UpdatedAt         string `db:"updated_at"`
}

// psql statement to update the players table to be new and improved
// ALTER TABLE players
// ADD COLUMN luck_status VARCHAR(255),
// ADD COLUMN alignment_override VARCHAR(255);

// ComplexPlayer is a player with a role
type ComplexPlayer struct {
	P Player `db:"players"`
	R Role   `db:"roles"`
}

// WARNING: May god have mercy on my soul for this abomination
type playerRoleJoin struct {
	PlayerID         int           `db:"player_id"`
	PlayerName       string        `db:"player_name"`
	PlayerGameID     string        `db:"player_game_id"`
	PlayerRoleID     int           `db:"player_role_id"`
	PlayerAlive      bool          `db:"player_alive"`
	PlayerSeat       int           `db:"player_seat"`
	PlayerLuck       int           `db:"player_luck"`
	PlayerLuckMod    int           `db:"player_luck_modifier"`
	PlayerLuckStatus string        `db:"player_luck_status"`
	PlayerAlignment  string        `db:"player_alignment_override"`
	PlayerCreated    string        `db:"player_created_at"`
	PlayerUpdated    string        `db:"player_updated_at"`
	RoleID           int           `db:"role_id"`
	RoleName         string        `db:"role_name"`
	RoleAlignment    string        `db:"role_alignment"`
	RoleAbilityIDs   pq.Int32Array `db:"role_ability_ids"`
	RolePassiveIDs   pq.Int32Array `db:"role_passive_ids"`
}

type PlayerModel struct {
	DB *sqlx.DB
}

func (m *PlayerModel) GetByGameID(gameID string) ([]*Player, error) {
	players := []*Player{}
	err := m.DB.Select(&players, "SELECT * FROM players WHERE game_id = $1", gameID)
	if err != nil {
		return nil, err
	}
	return players, nil
}

func (m *PlayerModel) GetByID(id int) (*Player, error) {
	player := &Player{}
	err := m.DB.Get(player, "SELECT * FROM players WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (m *PlayerModel) GetByName(name string) (*Player, error) {
	player := &Player{}
	err := m.DB.Get(player, "SELECT * FROM players WHERE name ILIKE $1", name)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (m *PlayerModel) GetPlayerNames(gameID string) ([]string, error) {
	names := []string{}
	err := m.DB.Select(&names, "SELECT name FROM players WHERE game_id = $1", gameID)
	if err != nil {
		return nil, err
	}
	return names, nil
}

func (m *PlayerModel) GetByGameIDAndName(gameID string, name string) (*Player, error) {
	player := &Player{}
	err := m.DB.Get(player, "SELECT * FROM players WHERE game_id = $1 AND name ILIKE $2", gameID, name)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (m *PlayerModel) Create(player *Player) error {
	_, err := m.DB.NamedExec(`INSERT INTO players 
    (name, game_id, role_id, alive, seat, luck, luck_modifier, luck_status, alignment_override) 
    VALUES (:name, :game_id, :role_id, :alive, :seat, :luck, :luck_modifier, :luck_status, :alignment_override)`, player)
	if err != nil {
		return err
	}
	return nil
}

func (m *PlayerModel) Update(player *Player) error {
	// _, err := m.DB.NamedExec("UPDATE players SET name = :name, game_id = :game_id, role_id = :role_id, alive = :alive, seat = :seat, luck = :luck, luck_modifier = :luck_modifier WHERE id = :id", player)
	_, err := m.DB.Exec(`UPDATE players SET 
    name = $1, game_id = $2, role_id = $3, alive = $4, seat = $5, luck = $6, luck_modifier = $7, luck_status = $8, alignment_override = $9 WHERE id = $10`,
		player.Name, player.GameID, player.RoleID, player.Alive, player.Seat, player.Luck, player.LuckModifier, player.LuckStatus, player.AlignmentOverride, player.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *PlayerModel) UpdateProperty(id int, property string, value interface{}) error {
	_, err := m.DB.Exec("UPDATE players SET $1 = $2 WHERE id = $3", property, value, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *PlayerModel) UpdateSeat(id int, seat int) error {
	_, err := m.DB.Exec("UPDATE players SET seat = $1 WHERE id = $2", seat, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *PlayerModel) Delete(id int) error {
	_, err := m.DB.Exec("DELETE FROM players WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (m *PlayerModel) GetRole(roleID int) (*Role, error) {
	role := &Role{}
	err := m.DB.Get(role, "SELECT * FROM roles WHERE id = $1", roleID)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (m *PlayerModel) GetComplexByGameID(gameID string, name string) (*ComplexPlayer, error) {
	playerQuery := &playerRoleJoin{}
	query := `SELECT p.id AS 
  player_id, p.name AS player_name, p.game_id AS player_game_id, p.role_id AS player_role_id, 
  p.alive AS player_alive, p.seat AS player_seat, p.luck AS player_luck, p.luck_modifier AS player_luck_modifier, 
  p.luck_status AS player_luck_status, p.alignment_override AS player_alignment_override, 
  p.created_at AS player_created_at, p.updated_at AS player_updated_at, 
  r.id AS role_id, r.name AS role_name, r.alignment AS role_alignment, r.ability_ids AS role_ability_ids, r.passive_ids AS role_passive_ids FROM players p JOIN roles r ON p.role_id = r.id WHERE p.game_id = $1 AND p.name ILIKE $2`
	err := m.DB.Get(playerQuery, query, gameID, name)
	if err != nil {
		return nil, err
	}
	player := &Player{
		ID:                playerQuery.PlayerID,
		Name:              playerQuery.PlayerName,
		GameID:            playerQuery.PlayerGameID,
		RoleID:            playerQuery.PlayerRoleID,
		Alive:             playerQuery.PlayerAlive,
		Seat:              playerQuery.PlayerSeat,
		Luck:              playerQuery.PlayerLuck,
		LuckModifier:      playerQuery.PlayerLuckMod,
		LuckStatus:        playerQuery.PlayerLuckStatus,
		AlignmentOverride: playerQuery.PlayerAlignment,
		CreatedAt:         playerQuery.PlayerCreated,
		UpdatedAt:         playerQuery.PlayerUpdated,
	}
	role := &Role{
		ID:         playerQuery.RoleID,
		Name:       playerQuery.RoleName,
		Alignment:  playerQuery.RoleAlignment,
		AbilityIDs: playerQuery.RoleAbilityIDs,
		PassiveIDs: playerQuery.RolePassiveIDs,
	}
	return &ComplexPlayer{
		P: *player,
		R: *role,
	}, nil
}

func (m *PlayerModel) GetAllComplexByGameID(gameID string) ([]*ComplexPlayer, error) {
	playerQuery := []*playerRoleJoin{}
	players := []*ComplexPlayer{}
	query := `SELECT p.id AS player_id, p.name AS player_name, p.game_id AS player_game_id, p.role_id AS player_role_id, 
  p.alive AS player_alive, p.seat AS player_seat, p.luck AS player_luck, 
  p.luck_modifier AS player_luck_modifier, 
  p.luck_status AS player_luck_status, p.alignment_override AS player_alignment_override,
  p.created_at AS player_created_at, p.updated_at AS player_updated_at, r.id AS role_id, r.name AS role_name, r.alignment AS role_alignment, r.ability_ids AS role_ability_ids, r.passive_ids AS role_passive_ids FROM players p JOIN roles r ON p.role_id = r.id WHERE p.game_id = $1`
	err := m.DB.Select(&playerQuery, query, gameID)
	if err != nil {
		return nil, err
	}

	for _, p := range playerQuery {
		player := &Player{
		ID:                p.PlayerID,
		Name:              p.PlayerName,
		GameID:            p.PlayerGameID,
		RoleID:            p.PlayerRoleID,
		Alive:             p.PlayerAlive,
		Seat:              p.PlayerSeat,
		Luck:              p.PlayerLuck,
		LuckModifier:      p.PlayerLuckMod,
		LuckStatus:        p.PlayerLuckStatus,
		AlignmentOverride: p.PlayerAlignment,
		CreatedAt:         p.PlayerCreated,
		UpdatedAt:         p.PlayerUpdated,
		}
		role := &Role{
			ID:         p.RoleID,
			Name:       p.RoleName,
			Alignment:  p.RoleAlignment,
			AbilityIDs: p.RoleAbilityIDs,
			PassiveIDs: p.RolePassiveIDs,
		}
		players = append(players, &ComplexPlayer{
			P: *player,
			R: *role,
		})
	}

	return players, nil
}

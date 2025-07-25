package hub_model

type Game struct {
	ID              string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name            string `gorm:"type:varchar(100);not null;unique" json:"name"`
	MaxPlayers      int    `gorm:"type:integer;not null;default:1" json:"max_players"`
	DurationMinutes int    `gorm:"type:integer;not null;default:30" json:"duration_minutes"`
	Description     string `gorm:"type:text" json:"description"`
	Status          string `gorm:"type:varchar(20);not null;default:'not_started';index" json:"status"`
}

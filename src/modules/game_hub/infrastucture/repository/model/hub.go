package hub_model

import "time"

type Hub struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	GameID    string    `gorm:"type:uuid;not null;index" json:"game_id"`
	HostID    string    `gorm:"type:uuid;not null;index" json:"host_id"`
	Status    string    `gorm:"type:varchar(20);not null;default:'not_started';index" json:"status"`
	StartTime time.Time `gorm:"type:timestamp;not null" json:"start_time"`
	EndTime   time.Time `gorm:"type:timestamp;default:null" json:"end_time"`
	EntryFee  float64   `gorm:"type:decimal(10,2);not null;default:0.00" json:"entry_fee"`
	Currency  string    `gorm:"type:varchar(3);not null;default:'USD'" json:"currency"`

	// Relations
	Game *Game `gorm:"foreignKey:GameID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"game,omitempty"`
}

package models

import "time"

type Playlist struct {
	ID        int `sql:"AUTO_INCREMENT"`
	Name      string
	CreatorIP string
	CreatedAt time.Time
	Songs     []Song `gorm:"many2many:playlist_song;"`
}

func NewPlaylist(name, ip string, createdAt time.Time) *Playlist {
	return &Playlist{
		Name:      name,
		CreatorIP: ip,
		CreatedAt: createdAt,
	}
}

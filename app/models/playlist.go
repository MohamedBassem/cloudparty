package models

import "time"

type Playlist struct {
	ID        int `sql:"AUTO_INCREMENT"`
	Name      string
	CreatorIP string
	CreatedAt time.Time
}

func NewPlaylist(name, ip string, createdAt time.Time) *Playlist {
	return &Playlist{
		Name:      name,
		CreatorIP: ip,
		CreatedAt: createdAt,
	}
}

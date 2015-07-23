package models

type Song struct {
	ID   int `sql:"AUTO_INCREMENT"`
	Name string
	Url  string
	Type string
}

func NewSong(name, url, t string) *Song {
	return &Song{
		Name: name,
		Url:  url,
		Type: t,
	}
}

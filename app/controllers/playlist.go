package controllers

import (
	"strings"
	"time"

	"golang.org/x/net/websocket"

	"github.com/MohamedBassem/CloudParty/app/models"
	"github.com/MohamedBassem/CloudParty/app/utils"
	"github.com/revel/revel"
)

type PlaylistController struct {
	*revel.Controller
}

func (c PlaylistController) Get(playlistId string) revel.Result {
	var playlist models.Playlist
	if DB.Where(&models.Playlist{Name: playlistId}).First(&playlist).RecordNotFound() {
		return c.NotFound("Playlist not found.")
	}
	var songs []models.Song
	DB.Model(&playlist).Related(&songs, "Songs")

	return c.Render(playlist, songs)
}

func (c PlaylistController) NewPlaylist() revel.Result {
	newId := utils.RandomString(20)
	playlist := models.NewPlaylist(newId, strings.Split(c.Request.RemoteAddr, ":")[0], time.Now())
	DB.Create(&playlist)
	return c.RenderJson(struct{ PlaylistID string }{newId})
}

func (c PlaylistController) NewSong(playlistId string) revel.Result {
	var playlist models.Playlist
	if DB.Where(&models.Playlist{Name: playlistId}).First(&playlist).RecordNotFound() {
		return c.NotFound("Playlist not found.")
	}

	songUrl := c.Params.Form.Get("url")
	var song models.Song
	if DB.Where(&models.Song{Url: songUrl}).First(&song).RecordNotFound() {
		song = *models.NewSong("", songUrl, "")
	}
	DB.Model(&playlist).Association("Songs").Append(&song)
	MainChan <- Message{Playlist: playlistId, Message: "ADD " + songUrl}

	c.Response.Status = 201
	return c.RenderText("")
}

func (c PlaylistController) Subscribe(playlistId string, ws *websocket.Conn) revel.Result {
	var playlist models.Playlist
	if DB.Where(&models.Playlist{Name: playlistId}).First(&playlist).RecordNotFound() {
		return c.NotFound("Playlist not found.")
	}

	myChan := make(chan string, 1000)
	SubscribePlaylist(playlistId, myChan)
	defer UnsubscribePlaylist(playlistId, myChan)

	for {
		msg := <-myChan
		err := websocket.Message.Send(ws, msg)
		if err != nil {
			revel.INFO.Println(msg)
			revel.INFO.Println(err)
			return nil
		}
	}

	return c.RenderText("")
}

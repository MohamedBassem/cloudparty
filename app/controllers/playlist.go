package controllers

import (
	"golang.org/x/net/websocket"

	"github.com/MohamedBassem/cloudparty/app/models"
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
	return c.Render(playlist)
}

func sendSongs(c chan string, songs []models.Song) {
	for _, song := range songs {
		c <- models.AddCommand{SongUrl: song.Url}.String()
	}
}

func listenSocket(ws *websocket.Conn, receivedChan chan string, doneChan chan struct{}) {
	var msg string
	for {
		err := websocket.Message.Receive(ws, &msg)
		if err != nil || msg == "" {
			doneChan <- struct{}{}
			revel.INFO.Printf("%s SOCKET CLOSED", ws.Request().RemoteAddr)
			return
		}
		revel.INFO.Printf("GOT '%s' FROM %s", msg, ws.Request().RemoteAddr)
		receivedChan <- msg

		// Exit if we are done
		select {
		case <-doneChan:
			return
		default:
		}
	}
}

func handleAddCommand(playlist models.Playlist, c models.AddCommand) {
	var song models.Song
	if DB.Where(&models.Song{Url: c.SongUrl}).First(&song).RecordNotFound() {
		song = *models.NewSong("", c.SongUrl, "")
	}
	DB.Model(&playlist).Association("Songs").Append(&song)
	PublishPlaylist(playlist.Name, c)
}

func parseAndPublish(playlist models.Playlist, msg string) {
	command := models.ParseCommand(msg)
	switch c := command.(type) {
	case models.AddCommand:
		handleAddCommand(playlist, c)
	}
}

func (c PlaylistController) Subscribe(playlistId string, ws *websocket.Conn) revel.Result {
	var playlist models.Playlist
	if DB.Where(&models.Playlist{Name: playlistId}).First(&playlist).RecordNotFound() {
		return c.NotFound("Playlist not found.")
	}

	var songs []models.Song
	DB.Model(&playlist).Related(&songs, "Songs")

	sendChan := make(chan string, 1000)
	receivedChan := make(chan string)
	doneChan := make(chan struct{})

	SubscribePlaylist(playlistId, sendChan)
	go listenSocket(ws, receivedChan, doneChan)

	// On exit unsubscribe from all chans
	defer UnsubscribePlaylist(playlistId, sendChan)

	go sendSongs(sendChan, songs)

MAINLOOP:
	for {
		select {
		case msg := <-sendChan:
			err := websocket.Message.Send(ws, msg)
			if err != nil {
				revel.INFO.Println(msg)
				revel.INFO.Println(err)
				return nil
			}
		case msg := <-receivedChan:
			parseAndPublish(playlist, msg)
		case <-doneChan:
			break MAINLOOP
		}
	}

	return nil
}

package controllers

import (
	"strings"
	"time"

	"github.com/MohamedBassem/CloudParty/app"
	"github.com/MohamedBassem/CloudParty/app/models"
	"github.com/MohamedBassem/CloudParty/app/utils"
	"github.com/revel/revel"
)

type Playlist struct {
	*revel.Controller
}

func (c Playlist) Get(playlistId string) revel.Result {
	return c.Render(playlistId)
}

func (c Playlist) New() revel.Result {
	newId := utils.RandomString(20)
	playlist := models.NewPlaylist(newId, strings.Split(c.Request.RemoteAddr, ":")[0], time.Now())
	app.DB.Create(&playlist)
	return c.RenderJson(struct{ PlaylistID string }{newId})
}

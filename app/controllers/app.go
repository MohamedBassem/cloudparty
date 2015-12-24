package controllers

import (
	"strings"
	"time"

	"github.com/MohamedBassem/cloudparty/app/models"
	"github.com/MohamedBassem/cloudparty/app/utils"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	newId := utils.RandomString(20)
	playlist := models.NewPlaylist(newId, strings.Split(c.Request.RemoteAddr, ":")[0], time.Now())
	DB.Create(&playlist)
	return c.Redirect("/pl/%s", newId)
}

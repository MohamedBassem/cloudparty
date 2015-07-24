package controllers

import (
	"sync"

	"github.com/revel/revel"
)

var MainChan chan Message

var Sockets map[string][]chan string
var locks map[string]*sync.Mutex
var lockSockets *sync.Mutex

type Message struct {
	Playlist string
	Message  string
}

func withLock(playlist string, f func()) {
	if locks[playlist] == nil {
		lockSockets.Lock()
		locks[playlist] = new(sync.Mutex)
		lockSockets.Unlock()
	}
	locks[playlist].Lock()
	defer locks[playlist].Unlock()
	f()
}

func worker() {
	var msg Message
	for {
		msg = <-MainChan
		subscribers := Sockets[msg.Playlist]
		for _, subscriber := range subscribers {
			subscriber <- msg.Message
		}
	}
}

func init() {
	MainChan = make(chan Message, 10000)
	Sockets = make(map[string][]chan string)
	locks = make(map[string]*sync.Mutex)
	lockSockets = new(sync.Mutex)

	for i := 0; i < 2; i++ {
		go worker()
	}
}

func SubscribePlaylist(playlist string, channel chan string) {
	withLock(playlist, func() {
		Sockets[playlist] = append(Sockets[playlist], channel)
		revel.INFO.Println("User Connected to channel " + playlist)
	})
}

func UnsubscribePlaylist(playlist string, channel chan string) {
	for i, val := range Sockets[playlist] {
		if val == channel {
			withLock(playlist, func() {
				copy(Sockets[playlist][i:], Sockets[playlist][i+1:])
				Sockets[playlist][len(Sockets[playlist])-1] = nil
				Sockets[playlist] = Sockets[playlist][:len(Sockets[playlist])-1]
				revel.INFO.Println("User Disconnected from channel " + playlist)
			})
		}
	}
}

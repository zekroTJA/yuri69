package discord

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

var statuses = []string{
	"annoying sounds",
	"funny sounds",
	"stupid sounds",
	"loud sounds",
	"horny sounds",
	"edgy sounds",
	"meme sounds",
	"insulting sounds",
	"kinky sounds",
}

func rotateStatus(s *discordgo.Session, every time.Duration) {
	ticker := time.NewTicker(every)

	go func() {
		for {
			status := statuses[rand.Intn(len(statuses))]
			s.UpdateGameStatus(0, status)
			<-ticker.C
		}
	}()
}
